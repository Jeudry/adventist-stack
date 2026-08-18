import 'dart:convert';

import 'package:http/http.dart' as http;

import '../config/app_config.dart';

/// Un campo concreto que el servidor rechazó, con dónde estaba.
class FieldError {
  const FieldError({required this.message, required this.location});

  factory FieldError.fromJson(Map<String, dynamic> json) => FieldError(
        message: json['message'] as String? ?? '',
        location: json['location'] as String? ?? '',
      );

  final String message;

  /// Ruta al valor que falló, como `body.firstName` o `query.pageSize`.
  final String location;
}

/// Excepción lanzada cuando el API responde con un código de error.
///
/// El servidor responde RFC 7807: `title` resume el problema, `detail` lo
/// explica, y `errors` trae un campo por cada valor rechazado — que es lo que
/// un formulario necesita para marcar la casilla concreta.
class ApiException implements Exception {
  ApiException(this.statusCode, this.message, {this.fieldErrors = const []});

  factory ApiException.fromProblem(int statusCode, Map<String, dynamic> body) {
    final errors = (body['errors'] as List<dynamic>? ?? <dynamic>[])
        .cast<Map<String, dynamic>>()
        .map(FieldError.fromJson)
        .toList(growable: false);

    return ApiException(
      statusCode,
      body['detail'] as String? ?? body['title'] as String? ?? 'Error desconocido',
      fieldErrors: errors,
    );
  }

  final int statusCode;
  final String message;
  final List<FieldError> fieldErrors;

  @override
  String toString() => 'ApiException($statusCode): $message';
}

/// Cliente HTTP mínimo hacia el API Gateway.
///
/// Guarda el access token en memoria y lo adjunta como Bearer en cada request.
/// La persistencia del token (secure storage) se agregará cuando construyamos
/// la feature de autenticación.
class ApiClient {
  ApiClient({http.Client? httpClient, String? baseUrl})
      : _http = httpClient ?? http.Client(),
        _baseUrl = baseUrl ?? AppConfig.apiBaseUrl;

  final http.Client _http;
  final String _baseUrl;
  String? _accessToken;

  set accessToken(String? token) => _accessToken = token;

  Map<String, String> get _headers => {
        'Content-Type': 'application/json',
        if (_accessToken != null) 'Authorization': 'Bearer $_accessToken',
      };

  Future<Map<String, dynamic>> get(
    String path, {
    Map<String, dynamic>? query,
  }) async {
    final res = await _http.get(_uri(path, query), headers: _headers);
    return _decode(res);
  }

  Future<Map<String, dynamic>> post(String path, Object body) async {
    final res = await _http.post(
      _uri(path),
      headers: _headers,
      body: jsonEncode(body),
    );
    return _decode(res);
  }

  Future<Map<String, dynamic>> put(String path, Object body) async {
    final res = await _http.put(
      _uri(path),
      headers: _headers,
      body: jsonEncode(body),
    );
    return _decode(res);
  }

  Future<void> delete(String path) async {
    _decode(await _http.delete(_uri(path), headers: _headers));
  }

  Uri _uri(String path, [Map<String, dynamic>? query]) {
    final uri = Uri.parse('$_baseUrl$path');
    if (query == null || query.isEmpty) return uri;
    return uri.replace(
      queryParameters: {
        for (final entry in query.entries)
          if (entry.value != null) entry.key: '${entry.value}',
      },
    );
  }

  Map<String, dynamic> _decode(http.Response res) {
    final body = res.body.isEmpty
        ? <String, dynamic>{}
        : jsonDecode(res.body) as Map<String, dynamic>;

    if (res.statusCode >= 400) {
      throw ApiException.fromProblem(res.statusCode, body);
    }
    return body;
  }
}
