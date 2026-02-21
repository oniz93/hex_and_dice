import 'dart:convert';
import '../models/messages.dart';
import '../models/game_state.dart';
import 'ws_service.dart';

class ParsedMessage {
  final String type;
  final dynamic data;
  final int? seq;

  ParsedMessage(this.type, this.data, this.seq);
}

ParsedMessage? parseMessage(ServerMessage msg) {
  try {
    final data = msg.data;
    if (data == null) return ParsedMessage(msg.type, null, msg.seq);

    switch (msg.type) {
      case 'game_state':
        return ParsedMessage(msg.type, GameState.fromJson(data), msg.seq);
      case 'ack':
        return ParsedMessage(msg.type, AckData.fromJson(data), msg.seq);
      case 'nack':
        return ParsedMessage(msg.type, NackData.fromJson(data), msg.seq);
      case 'troop_moved':
        return ParsedMessage(msg.type, TroopMovedData.fromJson(data), msg.seq);
      case 'combat_result':
        return ParsedMessage(
          msg.type,
          CombatResultData.fromJson(data),
          msg.seq,
        );
      case 'troop_purchased':
        return ParsedMessage(
          msg.type,
          TroopPurchasedData.fromJson(data),
          msg.seq,
        );
      case 'troop_destroyed':
        return ParsedMessage(
          msg.type,
          TroopDestroyedData.fromJson(data),
          msg.seq,
        );
      case 'structure_attacked':
        return ParsedMessage(
          msg.type,
          StructureAttackedData.fromJson(data),
          msg.seq,
        );
      case 'structure_fires':
        return ParsedMessage(
          msg.type,
          StructureFiresData.fromJson(data),
          msg.seq,
        );
      case 'turn_start':
        return ParsedMessage(msg.type, TurnStartData.fromJson(data), msg.seq);
      case 'game_over':
        return ParsedMessage(msg.type, GameOverData.fromJson(data), msg.seq);
      case 'match_found':
        return ParsedMessage(msg.type, MatchFoundData.fromJson(data), msg.seq);
      case 'error':
        return ParsedMessage(msg.type, ErrorData.fromJson(data), msg.seq);
      case 'player_disconnected':
      case 'player_reconnected':
      case 'ping':
      case 'pong':
      case 'emote':
        return ParsedMessage(msg.type, data, msg.seq);
      default:
        print('Unknown message type: ${msg.type}');
        return ParsedMessage(msg.type, data, msg.seq);
    }
  } catch (e, st) {
    print('Failed to parse message of type ${msg.type}: $e');
    print(st);
    return null;
  }
}
