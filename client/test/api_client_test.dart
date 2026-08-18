import 'package:adventist_stack_client/core/api/api_client.dart';
import 'package:flutter_test/flutter_test.dart';

/// Los cuerpos son respuestas reales del backend, copiadas de correrlo.
void main() {
  group('ApiException.fromProblem', () {
    test('lee el detalle de un problema simple', () {
      final error = ApiException.fromProblem(404, const {
        'title': 'Not Found',
        'status': 404,
        'detail': 'sabbath school not found',
      });

      expect(error.statusCode, 404);
      expect(error.message, 'sabbath school not found');
      expect(error.fieldErrors, isEmpty);
    });

    test('separa los campos que el servidor rechazó', () {
      final error = ApiException.fromProblem(422, const {
        'title': 'Unprocessable Entity',
        'status': 422,
        'detail': 'validation failed',
        'errors': [
          {'message': 'expected number <= 100', 'location': 'query.pageSize'},
          {'message': 'expected required property firstName', 'location': 'body'},
        ],
      });

      expect(error.fieldErrors, hasLength(2));
      expect(error.fieldErrors.first.location, 'query.pageSize');
      expect(error.fieldErrors.first.message, contains('<= 100'));
    });

    test('no revienta si el cuerpo no trae nada reconocible', () {
      final error = ApiException.fromProblem(500, const {});

      expect(error.message, 'Error desconocido');
      expect(error.fieldErrors, isEmpty);
    });

    test('cae al title cuando no hay detail', () {
      final error = ApiException.fromProblem(403, const {'title': 'Forbidden'});

      expect(error.message, 'Forbidden');
    });
  });
}
