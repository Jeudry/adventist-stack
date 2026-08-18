import 'json.dart';

enum SabbathSchoolStatus {
  active('active'),
  inactive('inactive');

  const SabbathSchoolStatus(this.wire);
  final String wire;

  static SabbathSchoolStatus fromJson(Object? value) =>
      SabbathSchoolStatus.values.firstWhere(
        (status) => status.wire == value,
        orElse: () => SabbathSchoolStatus.active,
      );
}

class SabbathSchool {
  const SabbathSchool({
    this.id,
    required this.name,
    this.teacherId,
    this.location,
    this.targetMinAge,
    this.targetMaxAge,
    this.status = SabbathSchoolStatus.active,
    this.createdAt,
    this.updatedAt,
  });

  factory SabbathSchool.fromJson(Map<String, dynamic> json) => SabbathSchool(
        id: json['id'] as String?,
        name: json['name'] as String? ?? '',
        teacherId: json['teacherId'] as String?,
        location: json['location'] as String?,
        targetMinAge: json['targetMinAge'] as int?,
        targetMaxAge: json['targetMaxAge'] as int?,
        status: SabbathSchoolStatus.fromJson(json['status']),
        createdAt: instantFromJson(json['createdAt']),
        updatedAt: instantFromJson(json['updatedAt']),
      );

  final String? id;
  final String name;
  final String? teacherId;
  final String? location;
  final int? targetMinAge;
  final int? targetMaxAge;
  final SabbathSchoolStatus status;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Map<String, dynamic> toJson() => withoutNulls({
        'name': name,
        'teacherId': teacherId,
        'location': location,
        'targetMinAge': targetMinAge,
        'targetMaxAge': targetMaxAge,
        'status': status.wire,
      });
}
