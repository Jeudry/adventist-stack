/// Configuración de la app inyectada en tiempo de compilación con --dart-define.
///
/// Ejemplo:
///   flutter run --dart-define=API_BASE_URL=http://localhost:8080
class AppConfig {
  /// URL base del API Gateway.
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'http://localhost:8080',
  );
}
