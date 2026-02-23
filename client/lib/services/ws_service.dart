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

  String? _pendingRoomId;
  String? _pendingToken;

  WsService({required this.baseUrl});

  Stream<WsConnectionState> get connectionState => _stateController.stream;
  Stream<ServerMessage> get messages => _messageController.stream;

  Future<void> connect(String token) async {
    if (_connectionState == WsConnectionState.connected ||
        _connectionState == WsConnectionState.connecting) {
      return;
    }

    _setConnectionState(WsConnectionState.connecting);
    final completer = Completer<void>();

    try {
      final wsUrl = baseUrl.replaceFirst('http', 'ws') + '/ws?token=$token';
      print('WsService: Connecting to $wsUrl');
      _channel = WebSocketChannel.connect(Uri.parse(wsUrl));

      _channel!.stream.listen(
        (data) {
          print('WsService: Received data: $data');
          if (!completer.isCompleted) {
            _setConnectionState(WsConnectionState.connected);
            completer.complete();
            print('WsService: Connection established (first message received)');
          }
          _handleMessage(data);
        },
        onDone: () {
          print('WsService: Stream done');
          if (!completer.isCompleted) {
            completer.completeError('Connection closed before ready');
          }
          _handleDisconnect();
        },
        onError: (error) {
          print('WsService: Stream error: $error');
          if (!completer.isCompleted) {
            completer.completeError(error);
          }
          _handleDisconnect();
        },
      );

      await completer.future;
    } catch (e) {
      print('WsService: Connect exception: $e');
      if (!completer.isCompleted) {
        completer.completeError(e);
      }
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

        if (msg.type == 'connected') {
          if (_pendingRoomId != null) {
            _sendImmediate('join_game', {'room_id': _pendingRoomId!});
            _pendingRoomId = null;
          }
          return;
        }

        if (msg.type == 'nack') {
          final nackData = msg.data as Map<String, dynamic>;
          if (nackData['action_type'] == 'join_game') {
            final roomId = nackData['data']?['room_id'] ?? _pendingRoomId;
            if (roomId != null) {
              print('WsService: join_game nacked, retrying in 500ms...');
              Future.delayed(const Duration(milliseconds: 500), () {
                _sendImmediate('join_game', {'room_id': roomId});
              });
              return;
            }
          }
        }

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

  void _sendImmediate(String type, Map<String, dynamic> data) {
    if (_channel == null) return;
    final msg = {'type': type, 'data': data};
    final encoded = jsonEncode(msg);
    print('WsService: Sending immediate: $encoded');
    _channel!.sink.add(encoded);
  }

  void _send(String type, Map<String, dynamic> data, {bool useSeq = false}) {
    print(
        'WsService: _send called for type $type, connectionState: $_connectionState, channel: ${_channel != null}');
    if (_connectionState != WsConnectionState.connected || _channel == null) {
      print('WsService: Dropping message $type because not connected.');
      return;
    }

    final msg = {'type': type, 'data': data};
    if (useSeq) {
      msg['seq'] = _nextSeq();
    }

    final encoded = jsonEncode(msg);
    print('WsService: Sending message: $encoded');
    _channel!.sink.add(encoded);
  }

  void sendJoinGame(String roomId) {
    if (_connectionState != WsConnectionState.connected) {
      _pendingRoomId = roomId;
      print(
          'WsService: Queued join_game for room $roomId, will send on connect');
      return;
    }
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
          'target_q': target.q.toInt(),
          'target_r': target.r.toInt(),
          'target_s': target.s.toInt(),
        },
        useSeq: true);
  }

  void sendAttack(String unitId, CubeCoord target) {
    _send(
        'attack',
        {
          'unit_id': unitId,
          'target_q': target.q.toInt(),
          'target_r': target.r.toInt(),
          'target_s': target.s.toInt(),
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
