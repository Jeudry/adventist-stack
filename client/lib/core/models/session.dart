enum Role {
  admin('admin'),
  member('member');

  const Role(this.wire);
  final String wire;

  static Role fromJson(Object? value) => Role.values.firstWhere(
        (role) => role.wire == value,
        orElse: () => Role.member,
      );
}

class User {
  const User({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
  });

  factory User.fromJson(Map<String, dynamic> json) => User(
        id: json['id'] as String? ?? '',
        email: json['email'] as String? ?? '',
        name: json['name'] as String? ?? '',
        role: Role.fromJson(json['role']),
      );

  final String id;
  final String email;
  final String name;
  final Role role;
}

/// Lo que devuelven `/auth/register` y `/auth/login`.
class Session {
  const Session({
    required this.user,
    required this.accessToken,
    required this.refreshToken,
  });

  factory Session.fromJson(Map<String, dynamic> json) => Session(
        user: User.fromJson(json['user'] as Map<String, dynamic>? ?? const {}),
        accessToken: json['accessToken'] as String? ?? '',
        refreshToken: json['refreshToken'] as String? ?? '',
      );

  final User user;
  final String accessToken;
  final String refreshToken;
}
