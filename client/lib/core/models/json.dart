/// Conversiones entre los formatos que usa el API y los tipos de Dart.
///
/// El API distingue dos cosas que en Dart son la misma clase: un instante
/// (`createdAt`, en RFC3339) y una fecha sin hora (`birthDate`, YYYY-MM-DD).
/// Serializar una fecha como instante le agrega una hora que el servidor
/// rechaza, así que cada una tiene su par de funciones.
library;

DateTime? dateFromJson(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}

String? dateToJson(DateTime? value) {
  if (value == null) return null;
  final month = value.month.toString().padLeft(2, '0');
  final day = value.day.toString().padLeft(2, '0');
  return '${value.year.toString().padLeft(4, '0')}-$month-$day';
}

DateTime? instantFromJson(Object? value) {
  if (value is! String || value.isEmpty) return null;
  return DateTime.tryParse(value);
}

/// Descarta las claves nulas para no mandar `null` donde el API espera que el
/// campo simplemente no venga.
Map<String, dynamic> withoutNulls(Map<String, dynamic> json) {
  return Map<String, dynamic>.fromEntries(
    json.entries.where((entry) => entry.value != null),
  );
}
