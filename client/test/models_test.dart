import 'package:adventist_stack_client/core/models/member.dart';
import 'package:adventist_stack_client/core/models/page.dart';
import 'package:adventist_stack_client/core/models/prayer.dart';
import 'package:adventist_stack_client/core/models/sabbath_school.dart';
import 'package:adventist_stack_client/core/models/session.dart';
import 'package:flutter_test/flutter_test.dart';

/// Los cuerpos de este archivo están copiados de `gateway/api/openapi.yaml`.
/// Si el API cambia un nombre y aquí no, estos tests son los que avisan.
void main() {
  group('Member', () {
    const body = {
      'id': '6f1c1e6e-0b6f-4b0e-9b0e-000000000001',
      'createdAt': '2026-08-07T10:00:00Z',
      'updatedAt': '2026-08-07T11:30:00Z',
      'firstName': 'Leandro',
      'lastName': 'Jimenez',
      'email': 'leandro@example.com',
      'gender': 'M',
      'birthDate': '1998-03-12',
      'status': 'visitor',
    };

    test('lee la respuesta del API', () {
      final member = Member.fromJson(body);

      expect(member.firstName, 'Leandro');
      expect(member.gender, Gender.male);
      expect(member.status, MemberStatus.visitor);
      expect(member.birthDate, DateTime(1998, 3, 12));
      expect(member.createdAt?.toUtc().hour, 10);
    });

    test('manda la fecha sin hora, que es lo que el API acepta', () {
      final json = Member.fromJson(body).toJson();

      expect(json['birthDate'], '1998-03-12');
      expect(json['status'], 'visitor');
      expect(json['gender'], 'M');
    });

    test('no manda campos vacíos ni los que pone el servidor', () {
      final json = const Member(firstName: 'Ana', lastName: 'Pérez').toJson();

      expect(json.containsKey('email'), isFalse);
      expect(json.containsKey('id'), isFalse);
      expect(json.containsKey('createdAt'), isFalse);
    });

    test('un estado desconocido no revienta la app', () {
      expect(Member.fromJson({'status': 'jubilado'}).status, MemberStatus.active);
    });
  });

  group('Prayer', () {
    test('lee la respuesta del API', () {
      final prayer = Prayer.fromJson(const {
        'id': 'p1',
        'title': 'Por mi familia',
        'description': 'Que estén bien',
        'authorName': 'Leandro',
        'isAnonymous': true,
        'status': 'answered',
      });

      expect(prayer.status, PrayerStatus.answered);
      expect(prayer.isAnonymous, isTrue);
      expect(prayer.toJson()['authorName'], 'Leandro');
    });
  });

  group('SabbathSchool', () {
    test('lee la respuesta del API', () {
      final school = SabbathSchool.fromJson(const {
        'name': 'Jóvenes',
        'targetMinAge': 15,
        'targetMaxAge': 25,
        'status': 'inactive',
      });

      expect(school.targetMinAge, 15);
      expect(school.status, SabbathSchoolStatus.inactive);
    });
  });

  group('Session', () {
    test('lee los tokens con los nombres que devuelve auth', () {
      final session = Session.fromJson(const {
        'user': {'id': 'u1', 'email': 'a@b.co', 'name': 'Ana', 'role': 'admin'},
        'accessToken': 'access',
        'refreshToken': 'refresh',
      });

      expect(session.accessToken, 'access');
      expect(session.refreshToken, 'refresh');
      expect(session.user.role, Role.admin);
    });
  });

  group('Page', () {
    test('lee una página de miembros', () {
      final page = Page.fromJson(const {
        'items': [
          {'firstName': 'Ana', 'lastName': 'Pérez'},
          {'firstName': 'Luis', 'lastName': 'Gómez'},
        ],
        'total': 42,
        'page': 1,
        'pageSize': 2,
      }, Member.fromJson);

      expect(page.items, hasLength(2));
      expect(page.items.first.firstName, 'Ana');
      expect(page.total, 42);
      expect(page.pageSize, 2);
    });
  });
}
