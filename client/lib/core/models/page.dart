/// Una página de resultados tal como la devuelve el API.
class Page<T> {
  const Page({
    required this.items,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory Page.fromJson(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) itemFromJson,
  ) {
    final items = (json['items'] as List<dynamic>? ?? <dynamic>[])
        .cast<Map<String, dynamic>>()
        .map(itemFromJson)
        .toList(growable: false);

    return Page(
      items: items,
      total: json['total'] as int? ?? 0,
      page: json['page'] as int? ?? 1,
      pageSize: json['pageSize'] as int? ?? items.length,
    );
  }

  final List<T> items;
  final int total;
  final int page;
  final int pageSize;
}
