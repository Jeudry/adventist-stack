import 'json.dart';

enum PrayerStatus {
  pending('pending'),
  answered('answered'),
  archived('archived');

  const PrayerStatus(this.wire);
  final String wire;

  static PrayerStatus fromJson(Object? value) => PrayerStatus.values.firstWhere(
        (status) => status.wire == value,
        orElse: () => PrayerStatus.pending,
      );
}

class Prayer {
  const Prayer({
    this.id,
    required this.title,
    required this.description,
    required this.authorName,
    this.isAnonymous = false,
    this.status = PrayerStatus.pending,
    this.createdAt,
    this.updatedAt,
  });

  factory Prayer.fromJson(Map<String, dynamic> json) => Prayer(
        id: json['id'] as String?,
        title: json['title'] as String? ?? '',
        description: json['description'] as String? ?? '',
        authorName: json['authorName'] as String? ?? '',
        isAnonymous: json['isAnonymous'] as bool? ?? false,
        status: PrayerStatus.fromJson(json['status']),
        createdAt: instantFromJson(json['createdAt']),
        updatedAt: instantFromJson(json['updatedAt']),
      );

  final String? id;
  final String title;
  final String description;
  final String authorName;
  final bool isAnonymous;
  final PrayerStatus status;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  Map<String, dynamic> toJson() => withoutNulls({
        'title': title,
        'description': description,
        'authorName': authorName,
        'isAnonymous': isAnonymous,
        'status': status.wire,
      });
}
