import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../models/messages.dart';
import '../models/enums.dart';
import '../game/hex/cube_coord.dart';

enum WsConnectionState { disconnected, connecting, connected, reconnecting }

class ServerMessage {
  final String type;
  final int? seq;
  final dynamic data;

  ServerMessage({required this.type, this.seq, this.data});

  factory ServerMessage.fromJson(Map<String, dynamic> json) {
    return ServerMessage(
      type: json['type'] as String,
      seq: json['seq'] as int?,
      data: json['data'],
    );
  }
}

class WsService {
  final String baseUrl;
  WebSocketChannel? _channel;
  WsConnectionState _connectionState = WsConnectionState.disconnected;

  final _stateController = StreamController<WsConnectionState>.broadcast();
  final _messageController = StreamController<ServerMessage>.broadcast();

  int _seqCounter = 0;
  Timer? _pongTimer;

  WsService({required this.baseUrl});

  Stream<WsConnectionState> get connectionState => _stateController.stream;
  Stream<ServerMessage> get messages => _messageController.stream;

  Future<void> connect(String token) async {
    if (_connectionState == WsConnectionState.connected ||
        _connectionState == WsConnectionState.connecting) {
      return;
    }

    _setConnectionState(WsConnectionState.connecting);
    try {
      final wsUrl = baseUrl.replaceFirst('http', 'ws') + '/ws?token=$token';
      _channel = WebSocketChannel.connect(Uri.parse(wsUrl));

      await _channel!.ready;
      _setConnectionState(WsConnectionState.connected);

      _channel!.stream.listen(
        (data) {
          _handleMessage(data);
        },
        onDone: () {
          _handleDisconnect();
        },
        onError: (error) {
          _handleDisconnect();
        },
      );
    } catch (e) {
      _handleDisconnect();
    }
  }

  void disconnect() {
    _channel?.sink.close();
    _channel = null;
    _setConnectionState(WsConnectionState.disconnected);
    _pongTimer?.cancel();
  }

  void _setConnectionState(WsConnectionState state) {
    if (_connectionState != state) {
      _connectionState = state;
      _stateController.add(state);
    }
  }

  void _handleDisconnect() {
    _channel = null;
    _setConnectionState(WsConnectionState.disconnected);
    _pongTimer?.cancel();
  }

  void _handleMessage(dynamic data) {
    if (data is String) {
      try {
        final decoded = jsonDecode(data);
        final msg = ServerMessage.fromJson(decoded);

        if (msg.type == 'ping') {
          sendPong();
          return;
        }

        _messageController.add(msg);
      } catch (e) {
        print('Error decoding WebSocket message: $e');
      }
    }
  }

  int _nextSeq() => ++_seqCounter;

  void _send(String type, Map<String, dynamic> data, {bool useSeq = false}) {
    if (_connectionState != WsConnectionState.connected || _channel == null)
      return;

    final msg = {'type': type, 'data': data};
    if (useSeq) {
      msg['seq'] = _nextSeq();
    }

    _channel!.sink.add(jsonEncode(msg));
  }

  void sendJoinGame(String roomId) {
    _send('join_game', {'room_id': roomId});
  }

  void sendReconnect(String gameId, String token) {
    _send('reconnect', {'game_id': gameId, 'player_token': token});
  }

  void sendMove(String unitId, CubeCoord target) {
    _send(
        'move',
        {
          'unit_id': unitId,
          'target_q': target.q,
          'target_r': target.r,
          'target_s': target.s,
        },
        useSeq: true);
  }

  void sendAttack(String unitId, CubeCoord target) {
    _send(
        'attack',
        {
          'unit_id': unitId,
          'target_q': target.q,
          'target_r': target.r,
          'target_s': target.s,
        },
        useSeq: true);
  }

  void sendBuy(TroopType type, String structureId) {
    _send(
        'buy',
        {
          'unit_type': type.name
              .replaceAll(RegExp(r'(?<!^)(?=[A-Z])'), '_')
              .toLowerCase(),
          'structure_id': structureId,
        },
        useSeq: true);
  }

  void sendEndTurn() {
    _send('end_turn', {}, useSeq: true);
  }

  void sendEmote(String emoteId) {
    _send('emote', {'emote_id': emoteId});
  }

  void sendPong() {
    _send('pong', {});
  }
}
