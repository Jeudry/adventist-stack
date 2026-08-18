import 'json.dart';

enum MemberStatus {
  active('active'),
  inactive('inactive'),
  visitor('visitor');

  const MemberStatus(this.wire);
  final String wire;

  static MemberStatus fromJson(Object? value) => MemberStatus.values.firstWhere(
        (status) => status.wire == value,
        orElse: () => MemberStatus.active,
      );
}

enum Gender {
  male('M'),
  female('F');

  const Gender(this.wire);
  final String wire;

  static Gender? fromJson(Object? value) {
    for (final gender in Gender.values) {
      if (gender.wire == value) return gender;
    }
    return null;
  }
}

class Member {
  const Member({
    this.id,
    required this.firstName,
    required this.lastName,
    this.email,
    this.phone,
    this.gender,
    this.address,
    this.birthDate,
    this.baptismDate,
    this.status = MemberStatus.active,
    this.createdAt,
    this.updatedAt,
  });

  factory Member.fromJson(Map<String, dynamic> json) => Member(
        id: json['id'] as String?,
        firstName: json['firstName'] as String? ?? '',
        lastName: json['lastName'] as String? ?? '',
        email: json['email'] as String?,
        phone: json['phone'] as String?,
        gender: Gender.fromJson(json['gender']),
        address: json['address'] as String?,
        birthDate: dateFromJson(json['birthDate']),
        baptismDate: dateFromJson(json['baptismDate']),
        status: MemberStatus.fromJson(json['status']),
        createdAt: instantFromJson(json['createdAt']),
        updatedAt: instantFromJson(json['updatedAt']),
      );

  final String? id;
  final String firstName;
  final String lastName;
  final String? email;
  final String? phone;
  final Gender? gender;
  final String? address;
  final DateTime? birthDate;
  final DateTime? baptismDate;
  final MemberStatus status;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  /// Solo los campos que el API acepta al crear o actualizar: el id y las
  /// marcas de auditoría las pone el servidor.
  Map<String, dynamic> toJson() => withoutNulls({
        'firstName': firstName,
        'lastName': lastName,
        'email': email,
        'phone': phone,
        'gender': gender?.wire,
        'address': address,
        'birthDate': dateToJson(birthDate),
        'baptismDate': dateToJson(baptismDate),
        'status': status.wire,
      });
}
