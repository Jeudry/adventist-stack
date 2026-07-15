# Guía de Protocol Buffers (protobuf) + gRPC

Material de referencia del proyecto. Explica qué es un `.proto`, cómo se lee, y
las convenciones **STANDARD** (linter `buf`) que seguimos.

---

## 1. ¿Qué es un `.proto` y por qué es "la fuente de verdad"?

**Protocol Buffers (protobuf)** son dos cosas a la vez:

1. Un **lenguaje para describir datos y APIs** — un *IDL* (Interface Definition Language).
2. Un **formato binario** para enviar esos datos por la red: más chico y rápido que JSON.

**gRPC** es el sistema de llamadas remotas (RPC) que usa protobuf para definir sus
servicios. Vos escribís **un** archivo `.proto` y, con `protoc`, se **genera código
automáticamente**: structs, cliente y servidor. El mismo `.proto` puede generar código
para Go, Dart/Flutter, TypeScript, Python…

> Por eso el `.proto` es "la fuente de verdad": manda el `.proto`, el código Go es un
> **derivado**. Por eso vive en `gen/` y **no se edita a mano**: si querés cambiar la
> API, cambiás el `.proto` y regenerás.

---

## 2. Anatomía de un `.proto`

```proto
syntax = "proto3";              // versión del LENGUAJE protobuf (usamos proto3)
package members.v1;             // namespace lógico + versión de LA API
option go_package = "...";      // dónde/con qué nombre generar el código Go

enum MemberStatus { ... }       // conjunto cerrado de valores
service MemberService { ... }   // un conjunto de operaciones (RPCs)
message Member { ... }          // una estructura de datos (como un struct)
```

- `syntax` → versión del lenguaje protobuf (proto3). **No confundir** con la versión
  de tu API (`v1`), que es otra cosa (ver §7).
- `package` → agrupa e identifica; termina en la versión de la API (`members.v1`).
- `option go_package` → `"<ruta-import>;<nombre-paquete-go>"`.

---

## 3. El concepto que MÁS confunde: los números de campo

```proto
message Member {
  string id = 1;          // ⚠️ el "= 1" NO es un valor: es el NÚMERO DE CAMPO
  string first_name = 2;
}
```

Ese número es la **identidad del campo en el formato binario**. Al serializar, protobuf
**no manda el nombre** `"first_name"`, manda el número `2`. Reglas de oro:

- Cada campo tiene un **número único** dentro del mensaje.
- Una vez en producción, **NO cambies ni reutilices** un número: romperías la
  compatibilidad con datos ya guardados o mensajes viejos en vuelo. El **nombre** sí lo
  podés cambiar (es cosmético); el **número es sagrado**.
- Los números **1–15** ocupan 1 byte; **16+** ocupan 2. Poné los campos más frecuentes
  del 1 al 15.

---

## 4. Tipos de campo

Un campo puede ser: un **escalar**, un **enum**, otro **mensaje** (composición), o una
colección (`repeated`, `map`). Además hay modificadores como `optional` y `oneof`.

### 4.1 Escalares — la lista completa

**Texto y binario**
| Tipo | Go | Uso |
|---|---|---|
| `string` | `string` | texto UTF-8 (nombres, emails, fechas ISO) |
| `bytes`  | `[]byte` | datos binarios crudos (imágenes, hashes) |
| `bool`   | `bool`   | verdadero/falso |

**Números enteros** (acá está la parte que sorprende a todos):
| Tipo | Go | Cuándo usarlo |
|---|---|---|
| `int32` / `int64` | `int32` / `int64` | enteros normales. Eficiente para valores **positivos**; ineficiente para negativos |
| `uint32` / `uint64` | `uint32` / `uint64` | enteros **sin signo** (nunca negativos) |
| `sint32` / `sint64` | `int32` / `int64` | enteros que **suelen ser negativos** (usan "zigzag", ver abajo) |
| `fixed32` / `fixed64` | `uint32` / `uint64` | siempre 4/8 bytes; conviene si el número suele ser **grande** |
| `sfixed32` / `sfixed64` | `int32` / `int64` | versión con signo de fixed |

**Números decimales**
| Tipo | Go |
|---|---|
| `float`  | `float32` |
| `double` | `float64` |

> **¿Por qué tantos enteros?** Por el formato binario. Los `int32/int64` se guardan como
> *varint*: un número chico ocupa pocos bytes. Pero un **negativo** en `int32` ocupa
> siempre 10 bytes (por cómo se representa). Si un campo puede ser negativo seguido
> (ej. una diferencia, una coordenada), usás `sint32` (codificación *zigzag*: mapea
> -1→1, 1→2, -2→3… para que los chicos, positivos o negativos, ocupen poco). Y si el
> número casi siempre es **grande** (ej. un hash, un id enorme), `fixed32/64` es mejor
> porque el varint ahí ya no ahorra. **Para el 95% de los casos, `int32`/`int64` está
> perfecto** — pero ahora sabés que las otras existen y por qué.

**Valores cero (defaults) por tipo:** `string`→`""`, `bytes`→`[]`, `bool`→`false`,
numéricos→`0`, enum→su valor 0, mensaje→`nil` (no seteado).

### 4.2 El "problema de la presencia" y `optional`

En proto3, un escalar **sin `optional`** no tiene *presencia*: si no lo mandás, llega con
su valor cero, y **no podés distinguir "no lo mandé" de "lo mandé en cero/vacío"**.

```proto
string email = 4;            // "" puede significar "vacío" O "no vino" — ambiguo
optional string email = 4;   // ahora SÍ se distingue
```

Con `optional`, protobuf **recuerda si el campo vino**. En Go se genera como **puntero**:

```go
m.Email          // *string
m.GetEmail()     // string ("" si es nil) — helper cómodo para leer
// nil            → "no vino"
// &""            → "vino, y es vacío"
```

**¿Cuándo lo necesitás de verdad?** En *updates parciales*:
- `email = nil` → "no toques el email"
- `email = &""` → "borrá el email (dejalo vacío)"

Sin `optional` no podrías ofrecer esa diferencia. (Por eso en `UpdateMemberRequest`
pusimos los campos como `optional`.)

> **Dato:** los campos de tipo **mensaje** (no escalar) SIEMPRE tienen presencia, incluso
> sin `optional` — en Go son punteros y `nil` = "no seteado". La presencia "faltaba"
> solo en los escalares; por eso `optional` es más relevante en ellos.

### 4.3 `repeated` — listas

```proto
repeated Member items = 1;   // -> []*Member en Go
repeated string tags = 2;    // -> []string
```
- Es una lista **ordenada**; puede estar vacía (nunca es `nil` conceptualmente, es `[]`).
- Para escalares numéricos, protobuf usa *packed encoding* (los empaqueta juntos) →
  más compacto. No tenés que hacer nada, es automático.

### 4.4 `map` — diccionarios

```proto
map<string, string> variables = 3;   // -> map[string]string en Go
```
La clave es un escalar (típico `string`); el valor puede ser escalar, enum o mensaje.
Ya lo usamos en `notifications` (`SendEmailRequest.variables`). Ojo: los `map` **no**
garantizan orden y **no** pueden ser `repeated`.

### 4.5 `enum` — conjunto cerrado de valores

```proto
enum MemberStatus {
  MEMBER_STATUS_UNSPECIFIED = 0;  // el 0 SIEMPRE es "sin especificar"
  MEMBER_STATUS_ACTIVE = 1;
  MEMBER_STATUS_INACTIVE = 2;
  MEMBER_STATUS_VISITOR = 3;
}
```
Mejor que un `string` "active"/"inactive": el compilador te limita a valores válidos.
Detalles importantes:
- **El valor 0 es obligatorio y por convención es `_UNSPECIFIED`** (es el default; sirve
  para detectar "no me mandaron nada").
- Los enums de proto3 son **"abiertos"**: si te llega un número que tu versión no conoce
  (porque el otro lado agregó `MEMBER_STATUS_DECEASED = 4`), **no explota**: lo conserva
  como entero. Esto es parte de la compatibilidad hacia adelante.
- En Go se genera como `type MemberStatus int32` + constantes
  (`membersv1.MemberStatus_MEMBER_STATUS_ACTIVE`).

### 4.6 Mensajes como campos (composición)

Un campo puede ser **otro mensaje** — así componés estructuras:
```proto
message Session {
  User user = 1;              // User es otro message -> *User en Go
  string access_token = 2;
}
```
Esto es exactamente lo que hicimos con `Session` dentro de `RegisterResponse`. Los
campos-mensaje son punteros en Go (`nil` = no seteado).

### 4.7 `oneof` — "uno de varios" (unión)

Cuando un campo puede ser **una cosa O la otra, pero nunca dos a la vez**:
```proto
message Contact {
  oneof method {
    string email = 1;
    string phone = 2;
  }
}
```

**Comportamiento clave:** los campos de un `oneof` **comparten espacio**. Si seteás
`email`, se limpia automáticamente `phone`, y viceversa. A lo sumo uno está presente.
Además, **siempre podés saber cuál** está seteado (tienen presencia por diseño), cosa que
con escalares sueltos no pasaba.

**Reglas:**
- Los números de campo siguen siendo únicos (1, 2, …), como en cualquier mensaje.
- Adentro de un `oneof` **no** podés poner `repeated` ni `map` directamente (si lo
  necesitás, envolvés eso en un mensaje aparte y ponés el mensaje en el `oneof`).
- No hace falta `optional` adentro: ya tienen presencia.

**Cómo se ve en Go** (esto es lo que confunde la primera vez): cada opción se vuelve un
tipo envoltorio, y el campo es una interfaz sobre la que hacés `switch`:
```go
// Setear:
c := &Contact{Method: &members.Contact_Email{Email: "ana@iglesia.org"}}

// Leer:
switch m := c.Method.(type) {
case *members.Contact_Email:
    fmt.Println("email:", m.Email)
case *members.Contact_Phone:
    fmt.Println("teléfono:", m.Phone)
case nil:
    fmt.Println("no vino ninguno")
}
```

**Cuándo lo usarías de verdad** (ejemplos del dominio):
- Una **notificación** que es email **o** SMS **o** push (cada canal con sus datos):
  ```proto
  message Notification {
    oneof channel {
      EmailPayload email = 1;
      SmsPayload   sms   = 2;
      PushPayload  push  = 3;
    }
  }
  ```
- Un **resultado** que es éxito **o** error.
- Una **búsqueda** que es por id **o** por email, pero no ambas.

`oneof` te da un tipo "unión" seguro: el compilador te obliga a manejar cada caso. No lo
usamos aún en el proyecto, pero es una herramienta muy práctica para payloads polimórficos.

### 4.8 Well-known types (tipos estándar de Google)

Protobuf trae mensajes ya definidos para casos comunes. Se **importan**:
```proto
import "google/protobuf/timestamp.proto";

message Event {
  google.protobuf.Timestamp starts_at = 1;   // fecha/hora estándar
}
```
Los más útiles: `Timestamp` (fecha/hora), `Duration` (duración), `Empty` (respuesta
vacía), `Struct` (JSON arbitrario). En `members` **usamos `Timestamp`** para las fechas
(`birth_date`, `baptism_date`, `created_at`, `updated_at`). En Go se generan como
`*timestamppb.Timestamp` y se convierten con:
```go
timestamppb.New(t)   // time.Time -> Timestamp
ts.AsTime()          // Timestamp -> time.Time
```

### 4.9 `service` y `rpc`

Hasta acá los `message` describen **datos**. Los `service` y `rpc` describen **acciones**:
qué operaciones puede pedir un cliente y qué manda/recibe en cada una.

#### ¿Qué es un RPC?

**RPC = Remote Procedure Call** (llamada a procedimiento remoto). La idea, que tiene
décadas, es simple y poderosa:

> Llamar a una función que **corre en OTRA máquina** como si fuera una función local.

Vos escribís algo que **parece** una llamada normal:
```go
res, err := memberClient.CreateMember(ctx, req)
```
…pero por debajo eso **no** ejecuta nada en tu proceso: empaqueta `req`, lo manda por la
red al otro servicio, ese lo ejecuta, y te devuelve la respuesta. Toda la plomería de red
(serializar, enviar, esperar, recibir, deserializar) queda **escondida**. Vos "llamás una
función"; gRPC hace la magia de que corra allá.

Ese es el punto de RPC: **hacer que hablar con otro servicio se sienta como llamar a un
método**, en vez de armar URLs y parsear JSON a mano.

#### ¿Por qué protobuf usa `service`?

Un `service` es simplemente **un grupo de RPCs relacionados** — el "menú" de operaciones
que ese microservicio ofrece. Es tu **contrato/interfaz**:

```proto
service MemberService {                                  // el contrato de "miembros"
  rpc CreateMember(CreateMemberRequest) returns (CreateMemberResponse);
  rpc GetMember(GetMemberRequest) returns (GetMemberResponse);
  //   ^nombre    ^qué recibe          ^qué devuelve
}
```

Pensalo como una **interfaz** de programación, pero que cruza la red. De ese `service`,
`protoc` genera **dos lados**:

- **Server** → una interfaz Go (`MemberServiceServer`) que **tu** microservicio
  implementa (ahí ponés la lógica real).
- **Client** → un objeto (`MemberServiceClient`) con métodos que **otros** llaman (ej. el
  gateway) y que por debajo hacen la llamada de red.

Los dos lados salen del **mismo** `.proto`, así que **no pueden desincronizarse**: si
cambiás la firma, ambos lados se regeneran y el compilador te avisa. Ese es el gran valor.

#### RPC vs REST (dos formas de pensar la API)

| | REST/HTTP | RPC/gRPC |
|---|---|---|
| Pensás en… | **recursos** y verbos (`GET /members/1`) | **funciones** (`GetMember(id)`) |
| El contrato es… | convención + docs (OpenAPI aparte) | el `.proto` (código generado) |
| Formato | JSON (texto) | protobuf (binario) |
| Tipado | manual | fuerte, autogenerado |

No es que uno sea "mejor": en este proyecto usamos **ambos**. Los microservicios hablan
entre sí por **gRPC** (rápido, tipado), y el **gateway** expone **REST** hacia afuera
(cómodo para el navegador y Flutter). El gateway es el traductor REST↔gRPC.

#### ¿Qué agrega gRPC al concepto de RPC?

RPC es la **idea**; **gRPC** es la implementación de Google que usamos. Combina:
- **protobuf** para serializar los mensajes (compacto y tipado),
- **HTTP/2** como transporte (rápido, multiplexado, soporta streaming).

Por eso ves `google.golang.org/grpc` en el código: es quien mueve los bytes por debajo.

#### Los 4 tipos de RPC

Lo que usamos es el RPC **unario** (una request → una response), como una función normal.
Pero gRPC soporta **streaming** (por eso existe la palabra clave `stream`):

| Tipo | Firma | Ejemplo de uso |
|---|---|---|
| **Unario** | `rpc F(Req) returns (Res)` | CRUD normal (lo nuestro) |
| **Server streaming** | `rpc F(Req) returns (stream Res)` | el server manda muchas respuestas (feed, notificaciones en vivo) |
| **Client streaming** | `rpc F(stream Req) returns (Res)` | el cliente sube muchos mensajes (subir un archivo por chunks) |
| **Bidireccional** | `rpc F(stream Req) returns (stream Res)` | chat en tiempo real |

Para esta app, unario nos alcanza. El streaming lo dejamos en el radar.

---

## 5. Convenciones STANDARD (linter `buf`)

Usamos `buf` con el conjunto de reglas **STANDARD** (`buf.yaml` en la raíz). Estas son
las que más se notan:

| Regla | Qué exige | Ejemplo |
|---|---|---|
| **Service suffix** | los servicios terminan en `Service` | `MemberService` ✓ |
| **RPC request/response standard name** | el request/response de un RPC se llama `<Método>Request` / `<Método>Response` | `CreateMember` → `CreateMemberRequest` / `CreateMemberResponse` |
| **RPC request/response unique** | un mismo mensaje **no** puede ser response de varios RPCs | por eso `Login` y `Register` NO comparten un `AuthResponse` |
| **Package version suffix** | el package termina en versión | `members.v1` ✓ (ver §7) |
| **Enum value prefix** | los valores del enum llevan el prefijo del enum | `MEMBER_STATUS_ACTIVE` ✓ |
| **Enum zero value suffix** | el valor 0 termina en `_UNSPECIFIED` | `MEMBER_STATUS_UNSPECIFIED = 0` ✓ |

**¿Por qué "response único por RPC"?** Si mañana `Login` necesita devolver un campo que
`Register` no (ej. `is_first_login`), con respuestas compartidas ensuciarías a ambos.
Separadas, cada endpoint evoluciona solo. Por eso, aunque hoy devuelvan lo mismo,
envolvemos los datos comunes en un mensaje anidado (ej. `Session`) y cada RPC tiene su
propio `...Response` que lo contiene:

```proto
message Session {                 // datos compartidos, anidados
  User user = 1;
  string access_token = 2;
  string refresh_token = 3;
}
message RegisterResponse { Session session = 1; }
message LoginResponse    { Session session = 1; }
```

Correr el linter:
```bash
buf lint      # exit 0 = todo conforme
```

---

## 6. De `.proto` a código Go

```bash
make proto
```
Ejecuta `protoc` sobre **todos** los `.proto`:
- `--go_out` → genera los **structs** (`Member`, `CreateMemberRequest`, el enum…) → `*.pb.go`
- `--go-grpc_out` → genera el **cliente + interfaz del servidor** gRPC → `*_grpc.pb.go`

Qué mirar en el código generado:
- Un `message` se vuelve un `type X struct { ... }`.
- Un `optional string` se vuelve `*string` (puntero → presencia).
- Un `enum` se vuelve `type X int32` con constantes.
- En `*_grpc.pb.go` está `XServiceServer` (la interfaz que **tu** servicio implementa) y
  `XServiceClient` (lo que usa el gateway).

---

## 7. ¿Por qué los protos tienen versiones? (`members.v1`)

Fijate que el package y la carpeta llevan `v1`: `package members.v1;`,
`proto/members/v1/`. Esto **no** es la versión de protobuf (eso es `proto3`); es la
**versión de TU API**. Existe por un problema central de los sistemas distribuidos:

> El cliente y el servidor se despliegan por separado y **no siempre a la vez**. Durante
> un rato conviven una versión vieja y una nueva. La API tiene que aguantar eso sin
> romperse.

### Cambios compatibles vs. incompatibles

Muchos cambios son **compatibles hacia atrás** y NO necesitan versión nueva:
- **Agregar** un campo nuevo (con un número nuevo). Los clientes viejos simplemente lo
  ignoran; los nuevos lo usan.
- **Agregar** un RPC o un valor de enum nuevo.

Otros cambios son **incompatibles (breaking)** y romperían a los clientes existentes:
- **Borrar** o **renombrar** un campo, o **cambiar su tipo** (`string` → `int32`).
- **Reutilizar un número** de campo para otra cosa.
- **Quitar** un RPC o cambiar su firma.

### Acá entra la versión

Para un cambio incompatible **no editás `v1`** (romperías a todos los que ya lo usan):
creás **`members/v2`** con `package members.v2;`, en paralelo. Entonces:

- Los clientes viejos siguen hablando `v1`.
- Los nuevos migran a `v2` cuando pueden.
- Cuando **nadie** usa más `v1`, lo eliminás.

La versión en el package te da un **namespace independiente**: `members.v1.Member` y
`members.v2.Member` son tipos distintos que conviven sin chocar. Es un contrato con tus
consumidores: *"v1 no va a romperse; los cambios grandes van a v2"*.

### Relación con los números de campo

Los números de campo dan compatibilidad **dentro** de una versión (por eso nunca se
reutilizan). La **versión del package** es la herramienta para cuando el cambio es tan
grande que ni los números alcanzan: empezás un contrato nuevo desde cero.

### En una frase
- `proto3` = versión del **lenguaje** protobuf.
- `v1` = versión de **tu API**; la subís a `v2` solo ante cambios que romperían a
  quienes ya usan `v1`.
- Números de campo = compatibilidad **fina**, dentro de una versión.

---

## 8. Reglas de oro para recordar

1. El número de campo es sagrado: no se cambia ni se reutiliza.
2. Agregar es (casi siempre) seguro; borrar/renombrar/cambiar-tipo es breaking → `v2`.
3. `optional` cuando necesites distinguir "no vino" de "vino vacío".
4. `enum`: el `0` siempre es `_UNSPECIFIED`.
5. Cada RPC, su `Request` y `Response` propios (STANDARD).
6. El `.proto` manda; `gen/` no se edita a mano.
