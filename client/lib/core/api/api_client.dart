import 'dart:convert';

import 'package:http/http.dart' as http;

import '../config/app_config.dart';

/// Excepción lanzada cuando el API responde con un código de error.
class ApiException implements Exception {
  ApiException(this.statusCode, this.message);

  final int statusCode;
  final String message;

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

  Future<Map<String, dynamic>> get(String path) async {
    final res = await _http.get(_uri(path), headers: _headers);
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

  Uri _uri(String path) => Uri.parse('$_baseUrl$path');

  Map<String, dynamic> _decode(http.Response res) {
    final body = res.body.isEmpty
        ? <String, dynamic>{}
        : jsonDecode(res.body) as Map<String, dynamic>;

    if (res.statusCode >= 400) {
      final message = body['error'] as String? ?? 'Error desconocido';
      throw ApiException(res.statusCode, message);
    }
    return body;
  }
}
