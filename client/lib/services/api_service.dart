import 'package:dio/dio.dart';
import '../models/enums.dart';
import '../models/room.dart';

class Session {
  final String id;
  final String nickname;
  final String token;

  Session({required this.id, required this.nickname, required this.token});

  factory Session.fromJson(Map<String, dynamic> json) {
    return Session(
      id: json['player_id'] ?? json['id'],
      nickname: json['nickname'],
      token: json['token'],
    );
  }
}

class RoomResponse {
  final String roomId;
  final String roomCode;
  final RoomSettings settings;
  final String? hostNickname;

  RoomResponse({
    required this.roomId,
    required this.roomCode,
    required this.settings,
    this.hostNickname,
  });

  factory RoomResponse.fromJson(Map<String, dynamic> json) {
    return RoomResponse(
      roomId: json['room_id'] ?? '',
      roomCode: json['room_code'] ?? '',
      settings: RoomSettings.fromJson(json['settings']),
      hostNickname: json['host_nickname'],
    );
  }
}

class RoomStatusResponse {
  final String code;
  final RoomState state;
  final RoomSettings settings;
  final String hostNickname;
  final String? guestNickname;
  final String? gameId;

  RoomStatusResponse({
    required this.code,
    required this.state,
    required this.settings,
    required this.hostNickname,
    this.guestNickname,
    this.gameId,
  });

  factory RoomStatusResponse.fromJson(Map<String, dynamic> json) {
    return RoomStatusResponse(
      code: json['code'],
      state: RoomState.values.firstWhere(
        (e) => _$RoomStateEnumMap[e] == json['state'],
      ),
      settings: RoomSettings.fromJson(json['settings']),
      hostNickname: json['host_nickname'],
      guestNickname: json['guest_nickname'],
      gameId: json['game_id'],
    );
  }
}

class MatchmakingResult {
  final String status;
  final String? roomId;
  final String? roomCode;

  MatchmakingResult({required this.status, this.roomId, this.roomCode});

  factory MatchmakingResult.fromJson(Map<String, dynamic> json) {
    return MatchmakingResult(
      status: json['status'],
      roomId: json['room_id'],
      roomCode: json['room_code'],
    );
  }
}

class ApiService {
  final Dio _dio;
  String? _token;

  ApiService({required String baseUrl})
    : _dio = Dio(
        BaseOptions(
          baseUrl: baseUrl,
          connectTimeout: const Duration(seconds: 5),
          receiveTimeout: const Duration(seconds: 5),
        ),
      ) {
    _dio.interceptors.add(
      InterceptorsWrapper(
        onRequest: (options, handler) {
          if (_token != null) {
            options.headers['Authorization'] = 'Bearer $_token';
          }
          return handler.next(options);
        },
      ),
    );
  }

  void setToken(String token) {
    _token = token;
  }

  Future<Session> registerGuest(String nickname) async {
    final response = await _dio.post(
      '/api/v1/guest',
      data: {'nickname': nickname},
    );
    return Session.fromJson(response.data);
  }

  Future<RoomResponse> createRoom(RoomSettings settings) async {
    final response = await _dio.post(
      '/api/v1/rooms',
      data: {
        'map_size': _$MapSizeEnumMap[settings.mapSize],
        'turn_timer': settings.turnTimer,
        'turn_mode': _$TurnModeEnumMap[settings.turnMode],
      },
    );
    return RoomResponse.fromJson(response.data);
  }

  Future<RoomResponse> joinRoom(String code) async {
    final response = await _dio.post(
      '/api/v1/rooms/join',
      data: {'code': code},
    );
    return RoomResponse.fromJson(response.data);
  }

  Future<RoomStatusResponse> getRoomStatus(String code) async {
    final response = await _dio.get('/api/v1/rooms/$code');
    return RoomStatusResponse.fromJson(response.data);
  }

  Future<MatchmakingResult> joinMatchmaking() async {
    final response = await _dio.post('/api/v1/matchmaking/join');
    return MatchmakingResult.fromJson(response.data);
  }

  Future<void> leaveMatchmaking() async {
    await _dio.delete('/api/v1/matchmaking/leave');
  }

  Future<Map<String, dynamic>> getMatchmakingStatus() async {
    final response = await _dio.get('/api/v1/matchmaking/status');
    return response.data as Map<String, dynamic>;
  }
}

const _$RoomStateEnumMap = {
  RoomState.waitingForOpponent: 'waiting_for_opponent',
  RoomState.ready: 'ready',
  RoomState.gameInProgress: 'game_in_progress',
  RoomState.gameOver: 'game_over',
};

const _$MapSizeEnumMap = {
  MapSize.small: 'small',
  MapSize.medium: 'medium',
  MapSize.large: 'large',
};

const _$TurnModeEnumMap = {
  TurnMode.alternating: 'alternating',
  TurnMode.simultaneous: 'simultaneous',
};
