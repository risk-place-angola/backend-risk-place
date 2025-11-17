# Guia de Implementação Flutter - Usuários Anônimos
## Risk Place Angola - Mobile Integration

**Versão Backend**: 1.0.0  
**Target**: Flutter 3.x+  
**Data**: Novembro 17, 2025

---

## 📱 Visão Geral

Este documento descreve **exatamente** como implementar o sistema de usuários anônimos no app Flutter, incluindo:
- Como gerar e armazenar o device ID
- Como registrar o dispositivo no backend
- Como conectar ao WebSocket
- Como receber e processar notificações
- Estrutura exata dos dados recebidos

---

## 🎯 Fluxo Completo de Integração

```
┌─────────────────────────────────────────────────────────────────┐
│                      APP FLUTTER INICIA                         │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  1. Gerar/Recuperar Device ID │
         │     UUID v4 persistente       │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  2. Obter FCM Token           │
         │     Firebase Messaging        │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  3. Obter Localização GPS     │
         │     Geolocator                │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  4. POST /devices/register    │
         │     Envia device_id + token   │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  5. Conectar WebSocket        │
         │     Header: X-Device-ID       │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │  6. Receber Notificações      │
         │     new_alert | report_created│
         └───────────────────────────────┘
```

---

## 📦 Dependências Necessárias

Adicione ao `pubspec.yaml`:

```yaml
dependencies:
  flutter:
    sdk: flutter
  
  # UUID generation
  uuid: ^4.2.1
  
  # Local storage
  shared_preferences: ^2.2.2
  
  # Location
  geolocator: ^11.0.0
  
  # Firebase Cloud Messaging
  firebase_core: ^2.24.2
  firebase_messaging: ^14.7.9
  
  # WebSocket
  web_socket_channel: ^2.4.0
  
  # HTTP requests
  http: ^1.1.2
  
  # JSON serialization
  json_annotation: ^4.8.1
  
  # State management (opcional)
  provider: ^6.1.1
  # OU
  bloc: ^8.1.3

dev_dependencies:
  build_runner: ^2.4.7
  json_serializable: ^6.7.1
```

---

## 🔧 1. Setup Inicial - Device ID Manager

Crie `lib/services/device_id_manager.dart`:

```dart
import 'package:shared_preferences/shared_preferences.dart';
import 'package:uuid/uuid.dart';

class DeviceIdManager {
  static const String _deviceIdKey = 'device_id';
  static const Uuid _uuid = Uuid();
  
  String? _cachedDeviceId;
  
  /// Obtém ou cria um device_id único e persistente
  Future<String> getDeviceId() async {
    // Retorna do cache se já carregado
    if (_cachedDeviceId != null) {
      return _cachedDeviceId!;
    }
    
    final prefs = await SharedPreferences.getInstance();
    
    // Tenta recuperar device_id existente
    String? deviceId = prefs.getString(_deviceIdKey);
    
    if (deviceId == null || deviceId.isEmpty) {
      // Gera novo UUID v4
      deviceId = _uuid.v4();
      await prefs.setString(_deviceIdKey, deviceId);
      print('✅ Novo device_id gerado: $deviceId');
    } else {
      print('✅ Device_id recuperado: $deviceId');
    }
    
    _cachedDeviceId = deviceId;
    return deviceId;
  }
  
  /// Limpa o device_id (apenas para debug/testes)
  Future<void> clearDeviceId() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_deviceIdKey);
    _cachedDeviceId = null;
    print('🗑️ Device_id removido');
  }
}
```

---

## 📡 2. Modelos de Dados

### 2.1 Request Models

Crie `lib/models/device_register_request.dart`:

```dart
import 'package:json_annotation/json_annotation.dart';

part 'device_register_request.g.dart';

@JsonSerializable()
class DeviceRegisterRequest {
  @JsonKey(name: 'device_id')
  final String deviceId;
  
  @JsonKey(name: 'fcm_token')
  final String? fcmToken;
  
  @JsonKey(name: 'platform')
  final String platform; // 'ios' ou 'android'
  
  @JsonKey(name: 'model')
  final String? model;
  
  @JsonKey(name: 'language')
  final String language; // 'pt'
  
  @JsonKey(name: 'latitude')
  final double? latitude;
  
  @JsonKey(name: 'longitude')
  final double? longitude;
  
  @JsonKey(name: 'alert_radius_meters')
  final int alertRadiusMeters;
  
  DeviceRegisterRequest({
    required this.deviceId,
    this.fcmToken,
    required this.platform,
    this.model,
    this.language = 'pt',
    this.latitude,
    this.longitude,
    this.alertRadiusMeters = 1000,
  });
  
  factory DeviceRegisterRequest.fromJson(Map<String, dynamic> json) =>
      _$DeviceRegisterRequestFromJson(json);
  
  Map<String, dynamic> toJson() => _$DeviceRegisterRequestToJson(this);
}
```

Crie `lib/models/update_location_request.dart`:

```dart
import 'package:json_annotation/json_annotation.dart';

part 'update_location_request.g.dart';

@JsonSerializable()
class UpdateLocationRequest {
  @JsonKey(name: 'device_id')
  final String deviceId;
  
  @JsonKey(name: 'latitude')
  final double latitude;
  
  @JsonKey(name: 'longitude')
  final double longitude;
  
  UpdateLocationRequest({
    required this.deviceId,
    required this.latitude,
    required this.longitude,
  });
  
  factory UpdateLocationRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateLocationRequestFromJson(json);
  
  Map<String, dynamic> toJson() => _$UpdateLocationRequestToJson(this);
}
```

### 2.2 Response Models

Crie `lib/models/device_register_response.dart`:

```dart
import 'package:json_annotation/json_annotation.dart';

part 'device_register_response.g.dart';

@JsonSerializable()
class DeviceRegisterResponse {
  @JsonKey(name: 'device_id')
  final String deviceId;
  
  @JsonKey(name: 'fcm_token')
  final String? fcmToken;
  
  @JsonKey(name: 'platform')
  final String? platform;
  
  @JsonKey(name: 'latitude')
  final double? latitude;
  
  @JsonKey(name: 'longitude')
  final double? longitude;
  
  @JsonKey(name: 'alert_radius_meters')
  final int? alertRadiusMeters;
  
  @JsonKey(name: 'message')
  final String? message;
  
  DeviceRegisterResponse({
    required this.deviceId,
    this.fcmToken,
    this.platform,
    this.latitude,
    this.longitude,
    this.alertRadiusMeters,
    this.message,
  });
  
  factory DeviceRegisterResponse.fromJson(Map<String, dynamic> json) =>
      _$DeviceRegisterResponseFromJson(json);
  
  Map<String, dynamic> toJson() => _$DeviceRegisterResponseToJson(this);
}
```

### 2.3 WebSocket Event Models

Crie `lib/models/websocket_events.dart`:

```dart
import 'package:json_annotation/json_annotation.dart';

part 'websocket_events.g.dart';

// Envelope genérico para mensagens WebSocket
@JsonSerializable(genericArgumentFactories: true)
class WebSocketMessage<T> {
  final String event;
  final T data;
  
  WebSocketMessage({
    required this.event,
    required this.data,
  });
  
  factory WebSocketMessage.fromJson(
    Map<String, dynamic> json,
    T Function(Object? json) fromJsonT,
  ) => _$WebSocketMessageFromJson(json, fromJsonT);
  
  Map<String, dynamic> toJson(Object Function(T value) toJsonT) =>
      _$WebSocketMessageToJson(this, toJsonT);
}

// Alerta recebido via WebSocket
@JsonSerializable()
class AlertNotification {
  @JsonKey(name: 'alert_id')
  final String alertId;
  
  @JsonKey(name: 'message')
  final String message;
  
  @JsonKey(name: 'latitude')
  final double latitude;
  
  @JsonKey(name: 'longitude')
  final double longitude;
  
  @JsonKey(name: 'radius')
  final double radius;
  
  AlertNotification({
    required this.alertId,
    required this.message,
    required this.latitude,
    required this.longitude,
    required this.radius,
  });
  
  factory AlertNotification.fromJson(Map<String, dynamic> json) =>
      _$AlertNotificationFromJson(json);
  
  Map<String, dynamic> toJson() => _$AlertNotificationToJson(this);
}

// Report recebido via WebSocket
@JsonSerializable()
class ReportNotification {
  @JsonKey(name: 'report_id')
  final String reportId;
  
  @JsonKey(name: 'message')
  final String message;
  
  @JsonKey(name: 'latitude')
  final double latitude;
  
  @JsonKey(name: 'longitude')
  final double longitude;
  
  ReportNotification({
    required this.reportId,
    required this.message,
    required this.latitude,
    required this.longitude,
  });
  
  factory ReportNotification.fromJson(Map<String, dynamic> json) =>
      _$ReportNotificationFromJson(json);
  
  Map<String, dynamic> toJson() => _$ReportNotificationToJson(this);
}

// Confirmação de atualização de localização
@JsonSerializable()
class LocationUpdatedResponse {
  final String status;
  
  LocationUpdatedResponse({required this.status});
  
  factory LocationUpdatedResponse.fromJson(Map<String, dynamic> json) =>
      _$LocationUpdatedResponseFromJson(json);
  
  Map<String, dynamic> toJson() => _$LocationUpdatedResponseToJson(this);
}
```

**IMPORTANTE**: Após criar os modelos, execute:
```bash
flutter pub run build_runner build --delete-conflicting-outputs
```

---

## 🌐 3. API Service

Crie `lib/services/api_service.dart`:

```dart
import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/device_register_request.dart';
import '../models/device_register_response.dart';
import '../models/update_location_request.dart';

class ApiService {
  // ⚠️ ALTERAR PARA SEU ENDPOINT
  static const String baseUrl = 'http://localhost:8000';
  // Production: 'https://api.riskplace.com'
  
  final http.Client _client;
  
  ApiService({http.Client? client}) : _client = client ?? http.Client();
  
  /// POST /api/v1/devices/register
  /// Registra o dispositivo anônimo no backend
  Future<DeviceRegisterResponse> registerDevice(
    DeviceRegisterRequest request,
  ) async {
    final url = Uri.parse('$baseUrl/api/v1/devices/register');
    
    print('📡 Registrando dispositivo: ${request.deviceId}');
    
    try {
      final response = await _client.post(
        url,
        headers: {
          'Content-Type': 'application/json',
        },
        body: jsonEncode(request.toJson()),
      );
      
      print('📡 Status: ${response.statusCode}');
      print('📡 Response: ${response.body}');
      
      if (response.statusCode == 200) {
        final json = jsonDecode(response.body) as Map<String, dynamic>;
        return DeviceRegisterResponse.fromJson(json);
      } else {
        throw Exception(
          'Falha ao registrar dispositivo: ${response.statusCode} - ${response.body}',
        );
      }
    } catch (e) {
      print('❌ Erro ao registrar dispositivo: $e');
      rethrow;
    }
  }
  
  /// PUT /api/v1/devices/location
  /// Atualiza a localização do dispositivo
  Future<void> updateDeviceLocation(
    UpdateLocationRequest request,
  ) async {
    final url = Uri.parse('$baseUrl/api/v1/devices/location');
    
    try {
      final response = await _client.put(
        url,
        headers: {
          'Content-Type': 'application/json',
        },
        body: jsonEncode(request.toJson()),
      );
      
      if (response.statusCode != 200) {
        throw Exception(
          'Falha ao atualizar localização: ${response.statusCode}',
        );
      }
      
      print('✅ Localização atualizada');
    } catch (e) {
      print('❌ Erro ao atualizar localização: $e');
      rethrow;
    }
  }
}
```

---

## 🔌 4. WebSocket Service

Crie `lib/services/websocket_service.dart`:

```dart
import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../models/websocket_events.dart';

class WebSocketService {
  // ⚠️ ALTERAR PARA SEU ENDPOINT
  static const String wsUrl = 'ws://localhost:8000/ws/alerts';
  // Production: 'wss://api.riskplace.com/ws/alerts'
  
  WebSocketChannel? _channel;
  final String deviceId;
  
  // Streams para notificações
  final _alertController = StreamController<AlertNotification>.broadcast();
  final _reportController = StreamController<ReportNotification>.broadcast();
  final _connectionController = StreamController<bool>.broadcast();
  
  Stream<AlertNotification> get alertStream => _alertController.stream;
  Stream<ReportNotification> get reportStream => _reportController.stream;
  Stream<bool> get connectionStream => _connectionController.stream;
  
  bool _isConnected = false;
  Timer? _reconnectTimer;
  Timer? _heartbeatTimer;
  int _reconnectAttempts = 0;
  
  WebSocketService({required this.deviceId});
  
  /// Conecta ao WebSocket com device_id
  Future<void> connect() async {
    if (_isConnected) {
      print('⚠️ WebSocket já conectado');
      return;
    }
    
    try {
      print('🔌 Conectando WebSocket...');
      print('🔌 URL: $wsUrl');
      print('🔌 Device ID: $deviceId');
      
      // IMPORTANTE: O device_id vai no header HTTP antes do upgrade
      final uri = Uri.parse(wsUrl);
      
      _channel = WebSocketChannel.connect(
        uri,
        // Note: web_socket_channel não suporta headers customizados diretamente
        // Solução: Enviar device_id na primeira mensagem
      );
      
      // Enviar device_id logo após conectar
      _sendDeviceRegistration();
      
      _isConnected = true;
      _reconnectAttempts = 0;
      _connectionController.add(true);
      
      print('✅ WebSocket conectado');
      
      // Iniciar heartbeat
      _startHeartbeat();
      
      // Escutar mensagens
      _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnect,
        cancelOnError: false,
      );
      
    } catch (e) {
      print('❌ Erro ao conectar WebSocket: $e');
      _isConnected = false;
      _connectionController.add(false);
      _scheduleReconnect();
    }
  }
  
  /// Envia device_id logo após conectar
  void _sendDeviceRegistration() {
    final message = jsonEncode({
      'event': 'register',
      'device_id': deviceId,
    });
    _channel?.sink.add(message);
    print('📤 Device ID enviado');
  }
  
  /// Processa mensagens recebidas do WebSocket
  void _handleMessage(dynamic message) {
    try {
      print('📨 Mensagem recebida: $message');
      
      final json = jsonDecode(message as String) as Map<String, dynamic>;
      final event = json['event'] as String?;
      final data = json['data'];
      
      if (event == null || data == null) {
        print('⚠️ Mensagem inválida: evento ou data ausente');
        return;
      }
      
      switch (event) {
        case 'new_alert':
          final alert = AlertNotification.fromJson(
            data as Map<String, dynamic>,
          );
          _alertController.add(alert);
          print('🚨 Alerta recebido: ${alert.message}');
          break;
          
        case 'report_created':
          final report = ReportNotification.fromJson(
            data as Map<String, dynamic>,
          );
          _reportController.add(report);
          print('📍 Report recebido: ${report.message}');
          break;
          
        case 'location_updated':
          print('✅ Localização confirmada pelo servidor');
          break;
          
        case 'pong':
          // Resposta ao heartbeat
          break;
          
        default:
          print('⚠️ Evento desconhecido: $event');
      }
    } catch (e) {
      print('❌ Erro ao processar mensagem: $e');
    }
  }
  
  /// Envia atualização de localização via WebSocket
  void updateLocation(double latitude, double longitude) {
    if (!_isConnected) {
      print('⚠️ WebSocket não conectado');
      return;
    }
    
    final message = jsonEncode({
      'event': 'update_location',
      'data': {
        'latitude': latitude,
        'longitude': longitude,
      },
    });
    
    _channel?.sink.add(message);
    print('📍 Localização enviada: $latitude, $longitude');
  }
  
  /// Heartbeat para manter conexão ativa
  void _startHeartbeat() {
    _heartbeatTimer?.cancel();
    _heartbeatTimer = Timer.periodic(
      const Duration(seconds: 30),
      (_) {
        if (_isConnected) {
          _channel?.sink.add(jsonEncode({'event': 'ping'}));
        }
      },
    );
  }
  
  /// Tratamento de erros
  void _handleError(error) {
    print('❌ WebSocket erro: $error');
    _isConnected = false;
    _connectionController.add(false);
    _scheduleReconnect();
  }
  
  /// Tratamento de desconexão
  void _handleDisconnect() {
    print('🔌 WebSocket desconectado');
    _isConnected = false;
    _connectionController.add(false);
    _heartbeatTimer?.cancel();
    _scheduleReconnect();
  }
  
  /// Agenda reconexão automática
  void _scheduleReconnect() {
    if (_reconnectTimer != null && _reconnectTimer!.isActive) {
      return;
    }
    
    _reconnectAttempts++;
    
    // Exponential backoff: 2s, 4s, 8s, 16s, max 60s
    final delay = Duration(
      seconds: (2 * _reconnectAttempts).clamp(2, 60),
    );
    
    print('🔄 Reconectando em ${delay.inSeconds}s (tentativa $_reconnectAttempts)');
    
    _reconnectTimer = Timer(delay, () {
      if (!_isConnected) {
        connect();
      }
    });
  }
  
  /// Desconecta o WebSocket
  void disconnect() {
    print('🔌 Desconectando WebSocket');
    _reconnectTimer?.cancel();
    _heartbeatTimer?.cancel();
    _channel?.sink.close();
    _isConnected = false;
    _connectionController.add(false);
  }
  
  /// Dispose de recursos
  void dispose() {
    disconnect();
    _alertController.close();
    _reportController.close();
    _connectionController.close();
  }
}
```

---

## 📍 5. Location Service

Crie `lib/services/location_service.dart`:

```dart
import 'dart:async';
import 'package:geolocator/geolocator.dart';

class LocationService {
  Position? _lastPosition;
  StreamSubscription<Position>? _positionSubscription;
  
  /// Solicita permissões de localização
  Future<bool> requestPermission() async {
    bool serviceEnabled;
    LocationPermission permission;
    
    // Verifica se serviço de localização está habilitado
    serviceEnabled = await Geolocator.isLocationServiceEnabled();
    if (!serviceEnabled) {
      print('❌ Serviço de localização desabilitado');
      return false;
    }
    
    // Verifica permissão
    permission = await Geolocator.checkPermission();
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
      if (permission == LocationPermission.denied) {
        print('❌ Permissão de localização negada');
        return false;
      }
    }
    
    if (permission == LocationPermission.deniedForever) {
      print('❌ Permissão de localização negada permanentemente');
      return false;
    }
    
    print('✅ Permissão de localização concedida');
    return true;
  }
  
  /// Obtém posição atual
  Future<Position?> getCurrentPosition() async {
    try {
      _lastPosition = await Geolocator.getCurrentPosition(
        desiredAccuracy: LocationAccuracy.high,
      );
      return _lastPosition;
    } catch (e) {
      print('❌ Erro ao obter localização: $e');
      return null;
    }
  }
  
  /// Inicia tracking contínuo de localização
  Stream<Position> startTracking({
    Duration interval = const Duration(seconds: 30),
  }) {
    const locationSettings = LocationSettings(
      accuracy: LocationAccuracy.high,
      distanceFilter: 10, // Metros
    );
    
    return Geolocator.getPositionStream(
      locationSettings: locationSettings,
    );
  }
  
  /// Para tracking
  void stopTracking() {
    _positionSubscription?.cancel();
    _positionSubscription = null;
  }
  
  Position? get lastPosition => _lastPosition;
}
```

---

## 🎬 6. Anonymous User Manager (Orquestrador Principal)

Crie `lib/services/anonymous_user_manager.dart`:

```dart
import 'dart:async';
import 'dart:io';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:geolocator/geolocator.dart';
import 'device_id_manager.dart';
import 'api_service.dart';
import 'websocket_service.dart';
import 'location_service.dart';
import '../models/device_register_request.dart';
import '../models/update_location_request.dart';

/// Gerenciador principal para usuários anônimos
class AnonymousUserManager {
  final DeviceIdManager _deviceIdManager;
  final ApiService _apiService;
  final LocationService _locationService;
  
  WebSocketService? _websocketService;
  Timer? _locationUpdateTimer;
  StreamSubscription<Position>? _locationSubscription;
  
  String? _deviceId;
  String? _fcmToken;
  
  AnonymousUserManager({
    DeviceIdManager? deviceIdManager,
    ApiService? apiService,
    LocationService? locationService,
  })  : _deviceIdManager = deviceIdManager ?? DeviceIdManager(),
        _apiService = apiService ?? ApiService(),
        _locationService = locationService ?? LocationService();
  
  /// Inicializa todo o sistema de usuário anônimo
  Future<void> initialize() async {
    try {
      print('🚀 Inicializando sistema de usuário anônimo...');
      
      // 1. Obter Device ID
      _deviceId = await _deviceIdManager.getDeviceId();
      print('✅ Device ID: $_deviceId');
      
      // 2. Solicitar permissão de localização
      final hasPermission = await _locationService.requestPermission();
      if (!hasPermission) {
        print('⚠️ Sem permissão de localização - funcionalidade limitada');
      }
      
      // 3. Obter FCM Token
      try {
        _fcmToken = await FirebaseMessaging.instance.getToken();
        print('✅ FCM Token: $_fcmToken');
      } catch (e) {
        print('⚠️ Erro ao obter FCM token: $e');
      }
      
      // 4. Registrar dispositivo no backend
      await _registerDevice();
      
      // 5. Conectar WebSocket
      await _connectWebSocket();
      
      // 6. Iniciar tracking de localização
      _startLocationTracking();
      
      // 7. Configurar listeners FCM
      _setupFCMListeners();
      
      print('✅ Sistema de usuário anônimo inicializado com sucesso!');
      
    } catch (e) {
      print('❌ Erro ao inicializar: $e');
      rethrow;
    }
  }
  
  /// Registra o dispositivo no backend
  Future<void> _registerDevice() async {
    final position = await _locationService.getCurrentPosition();
    
    final request = DeviceRegisterRequest(
      deviceId: _deviceId!,
      fcmToken: _fcmToken,
      platform: Platform.isIOS ? 'ios' : 'android',
      model: Platform.isIOS ? 'iOS' : 'Android', // Pode usar device_info_plus para mais detalhes
      language: 'pt',
      latitude: position?.latitude,
      longitude: position?.longitude,
      alertRadiusMeters: 1000, // 1km
    );
    
    final response = await _apiService.registerDevice(request);
    print('✅ Dispositivo registrado: ${response.message}');
  }
  
  /// Conecta ao WebSocket
  Future<void> _connectWebSocket() async {
    _websocketService = WebSocketService(deviceId: _deviceId!);
    await _websocketService!.connect();
    
    // Escutar alertas
    _websocketService!.alertStream.listen((alert) {
      print('🚨 ALERTA: ${alert.message}');
      _handleAlert(alert);
    });
    
    // Escutar reports
    _websocketService!.reportStream.listen((report) {
      print('📍 REPORT: ${report.message}');
      _handleReport(report);
    });
    
    // Escutar status de conexão
    _websocketService!.connectionStream.listen((isConnected) {
      print('🔌 WebSocket: ${isConnected ? "Conectado" : "Desconectado"}');
    });
  }
  
  /// Inicia tracking de localização
  void _startLocationTracking() {
    // Atualizar a cada 30 segundos
    _locationUpdateTimer = Timer.periodic(
      const Duration(seconds: 30),
      (_) async {
        await _updateLocation();
      },
    );
    
    // Também escutar mudanças significativas
    _locationSubscription = _locationService.startTracking().listen(
      (position) {
        // Atualizar via WebSocket (tempo real)
        _websocketService?.updateLocation(
          position.latitude,
          position.longitude,
        );
      },
    );
  }
  
  /// Atualiza localização no backend
  Future<void> _updateLocation() async {
    final position = await _locationService.getCurrentPosition();
    if (position == null) return;
    
    try {
      // Atualizar via HTTP (persistência)
      final request = UpdateLocationRequest(
        deviceId: _deviceId!,
        latitude: position.latitude,
        longitude: position.longitude,
      );
      
      await _apiService.updateDeviceLocation(request);
      
      // Atualizar via WebSocket (tempo real)
      _websocketService?.updateLocation(
        position.latitude,
        position.longitude,
      );
      
    } catch (e) {
      print('❌ Erro ao atualizar localização: $e');
    }
  }
  
  /// Configura listeners do Firebase Cloud Messaging
  void _setupFCMListeners() {
    // Mensagem recebida quando app está em foreground
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      print('📩 Notificação em foreground:');
      print('   Título: ${message.notification?.title}');
      print('   Corpo: ${message.notification?.body}');
      print('   Data: ${message.data}');
      
      _showLocalNotification(message);
    });
    
    // App aberto via notificação
    FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
      print('📱 App aberto via notificação');
      _handleNotificationTap(message);
    });
  }
  
  /// Processa alerta recebido
  void _handleAlert(AlertNotification alert) {
    // TODO: Implementar lógica específica do app
    // - Mostrar notificação local
    // - Atualizar UI
    // - Adicionar marcador no mapa
    // - Tocar som de alerta
  }
  
  /// Processa report recebido
  void _handleReport(ReportNotification report) {
    // TODO: Implementar lógica específica do app
    // - Mostrar notificação local
    // - Atualizar UI
    // - Adicionar marcador no mapa
  }
  
  /// Mostra notificação local
  void _showLocalNotification(RemoteMessage message) {
    // TODO: Usar flutter_local_notifications
    print('🔔 Mostrar notificação: ${message.notification?.title}');
  }
  
  /// Trata tap em notificação
  void _handleNotificationTap(RemoteMessage message) {
    // TODO: Navegar para tela específica baseado em message.data
    print('👆 Notificação clicada: ${message.data}');
  }
  
  /// Limpa recursos
  void dispose() {
    _locationUpdateTimer?.cancel();
    _locationSubscription?.cancel();
    _locationService.stopTracking();
    _websocketService?.dispose();
  }
  
  // Getters
  String? get deviceId => _deviceId;
  WebSocketService? get websocketService => _websocketService;
  Stream<AlertNotification>? get alertStream => _websocketService?.alertStream;
  Stream<ReportNotification>? get reportStream => _websocketService?.reportStream;
}
```

---

## 🎯 7. Uso no App

### 7.1 Inicialização no main.dart

```dart
import 'package:flutter/material.dart';
import 'package:firebase_core/firebase_core.dart';
import 'services/anonymous_user_manager.dart';

// Singleton global
final anonymousUserManager = AnonymousUserManager();

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Inicializar Firebase
  await Firebase.initializeApp();
  
  // Inicializar sistema de usuário anônimo
  try {
    await anonymousUserManager.initialize();
  } catch (e) {
    print('Erro na inicialização: $e');
  }
  
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Risk Place Angola',
      home: HomePage(),
    );
  }
}
```

### 7.2 Exemplo de Widget que Escuta Notificações

```dart
import 'package:flutter/material.dart';
import '../main.dart';
import '../models/websocket_events.dart';

class AlertListenerWidget extends StatefulWidget {
  @override
  _AlertListenerWidgetState createState() => _AlertListenerWidgetState();
}

class _AlertListenerWidgetState extends State<AlertListenerWidget> {
  final List<AlertNotification> _alerts = [];
  final List<ReportNotification> _reports = [];
  
  @override
  void initState() {
    super.initState();
    
    // Escutar alertas
    anonymousUserManager.alertStream?.listen((alert) {
      setState(() {
        _alerts.insert(0, alert);
      });
      
      // Mostrar SnackBar
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('🚨 ${alert.message}'),
          backgroundColor: Colors.red,
          duration: Duration(seconds: 5),
          action: SnackBarAction(
            label: 'Ver',
            onPressed: () {
              // Navegar para mapa ou detalhes
            },
          ),
        ),
      );
    });
    
    // Escutar reports
    anonymousUserManager.reportStream?.listen((report) {
      setState(() {
        _reports.insert(0, report);
      });
      
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('📍 ${report.message}'),
          backgroundColor: Colors.orange,
        ),
      );
    });
  }
  
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Notificações')),
      body: ListView(
        children: [
          if (_alerts.isNotEmpty) ...[
            Padding(
              padding: EdgeInsets.all(16),
              child: Text(
                'Alertas',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
            ),
            ..._alerts.map((alert) => ListTile(
              leading: Icon(Icons.warning, color: Colors.red),
              title: Text(alert.message),
              subtitle: Text('${alert.latitude}, ${alert.longitude}'),
              trailing: Text('${alert.radius.toInt()}m'),
            )),
          ],
          if (_reports.isNotEmpty) ...[
            Padding(
              padding: EdgeInsets.all(16),
              child: Text(
                'Reports',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
            ),
            ..._reports.map((report) => ListTile(
              leading: Icon(Icons.location_on, color: Colors.orange),
              title: Text(report.message),
              subtitle: Text('${report.latitude}, ${report.longitude}'),
            )),
          ],
        ],
      ),
    );
  }
}
```

---

## 📋 8. Estrutura Exata dos Dados Recebidos

### 8.1 Alerta via WebSocket

```json
{
  "event": "new_alert",
  "data": {
    "alert_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "🚨 Assalto reportado na área - Zona de Maianga",
    "latitude": -8.8390,
    "longitude": 13.2345,
    "radius": 500.0
  }
}
```

### 8.2 Report via WebSocket

```json
{
  "event": "report_created",
  "data": {
    "report_id": "660f9511-f3ac-52e5-b827-557766551111",
    "message": "📍 Buraco grande na via - Avenida 4 de Fevereiro",
    "latitude": -8.8395,
    "longitude": 13.2348
  }
}
```

### 8.3 Push Notification via FCM

```json
{
  "notification": {
    "title": "🚨 Alerta de Risco",
    "body": "Assalto reportado na área"
  },
  "data": {
    "alert_id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "alert",
    "latitude": "-8.8390",
    "longitude": "13.2345",
    "radius": "500"
  }
}
```

---

## ✅ 9. Checklist de Implementação

- [ ] Adicionar dependências no `pubspec.yaml`
- [ ] Criar `DeviceIdManager` e gerar UUID persistente
- [ ] Criar modelos de dados com `json_serializable`
- [ ] Rodar `build_runner` para gerar `*.g.dart`
- [ ] Implementar `ApiService` com endpoints corretos
- [ ] Implementar `WebSocketService` com tratamento de eventos
- [ ] Implementar `LocationService` com permissões
- [ ] Criar `AnonymousUserManager` como orquestrador
- [ ] Inicializar no `main.dart`
- [ ] Configurar Firebase (google-services.json / GoogleService-Info.plist)
- [ ] Testar registro de dispositivo
- [ ] Testar conexão WebSocket
- [ ] Testar recebimento de notificações
- [ ] Implementar UI para exibir alertas/reports
- [ ] Adicionar notificações locais (flutter_local_notifications)
- [ ] Testar em background/foreground
- [ ] Implementar reconexão automática

---

## 🧪 10. Como Testar

### Teste 1: Registro de Dispositivo

1. Rode o app
2. Verifique logs:
   ```
   ✅ Device ID: 550e8400-...
   ✅ FCM Token: dQw4w9...
   ✅ Dispositivo registrado: Device registered successfully
   ```

### Teste 2: WebSocket

1. Verifique logs:
   ```
   🔌 Conectando WebSocket...
   📤 Device ID enviado
   ✅ WebSocket conectado
   ```

### Teste 3: Receber Alerta

Simule via backend ou curl:
```bash
# Criar alerta próximo da localização do dispositivo
curl -X POST http://localhost:8000/api/v1/alerts \
  -H "Authorization: Bearer <ERCE_JWT>" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Teste de alerta",
    "latitude": -8.8390,
    "longitude": 13.2345,
    "radius": 500
  }'
```

Deve aparecer no app:
```
📨 Mensagem recebida: {"event":"new_alert","data":{...}}
🚨 Alerta recebido: Teste de alerta
```

---

## 🚨 Troubleshooting

### Problema: WebSocket não conecta

**Solução**:
```dart
// Verificar URL (ws:// para desenvolvimento, wss:// para produção)
// Android: usar 10.0.2.2 em vez de localhost
static const String wsUrl = 'ws://10.0.2.2:8000/ws/alerts';
```

### Problema: FCM Token nulo

**Solução**:
- Verificar `google-services.json` (Android) / `GoogleService-Info.plist` (iOS)
- Verificar SHA-1 no Firebase Console
- Testar em dispositivo real

### Problema: Permissão de localização negada

**Solução**:
```dart
// Adicionar ao AndroidManifest.xml
<uses-permission android:name="android.permission.ACCESS_FINE_LOCATION" />
<uses-permission android:name="android.permission.ACCESS_COARSE_LOCATION" />

// iOS: info.plist
<key>NSLocationWhenInUseUsageDescription</key>
<string>Precisamos da sua localização para enviar alertas próximos</string>
```

---

## 📚 Referências

- [Backend README](./ANONYMOUS_USERS_README.md)
- [Guia Completo](./ANONYMOUS_USER_GUIDE.md)
- [Arquitetura](./diagram/ANONYMOUS_USERS_ARCHITECTURE.md)
- [Geolocator Docs](https://pub.dev/packages/geolocator)
- [Firebase Messaging Docs](https://pub.dev/packages/firebase_messaging)
- [WebSocket Channel Docs](https://pub.dev/packages/web_socket_channel)

---

**Versão**: 1.0.0  
**Autor**: Backend Team  
**Contato**: Para dúvidas, consulte a equipe de backend ou abra issue no repositório
