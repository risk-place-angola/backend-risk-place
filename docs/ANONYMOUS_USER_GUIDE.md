# Guia de Integração para Usuários Anônimos

**Versão**: 1.0.0  
**Data**: Novembro 17, 2025

---

## Visão Geral

O sistema Risk Place Angola agora suporta **usuários anônimos** (não autenticados) que podem receber notificações de alertas e reports sem necessidade de criar conta ou fazer login, similar ao funcionamento do Waze.

## Como Funciona

### Arquitetura de Usuários Anônimos

```
┌─────────────────────┐
│   Mobile App        │
│   (Sem Login)       │
└──────────┬──────────┘
           │
           │ 1. POST /api/v1/devices/register
           │    { device_id, fcm_token, location }
           ▼
┌─────────────────────┐
│   Backend API       │
│   Cria/Atualiza     │
│   AnonymousSession  │
└──────────┬──────────┘
           │
           │ 2. WebSocket Connect
           │    ws://host/ws/alerts
           │    Header: X-Device-ID: <unique_device_id>
           ▼
┌─────────────────────┐
│  WebSocket Hub      │
│  Registra Cliente   │
│  Anônimo            │
└──────────┬──────────┘
           │
           │ 3. Recebe Notificações
           │    - Alertas próximos
           │    - Reports próximos
           ▼
┌─────────────────────┐
│  Push Notifications │
│  (FCM - Fallback)   │
└─────────────────────┘
```

---

## Fluxo de Implementação Mobile

### 1. Gerar Device ID Único

O `device_id` deve ser **único por dispositivo** e persistente:

```dart
// Flutter exemplo
import 'package:uuid/uuid.dart';
import 'package:shared_preferences/shared_preferences.dart';

Future<String> getOrCreateDeviceId() async {
  final prefs = await SharedPreferences.getInstance();
  String? deviceId = prefs.getString('device_id');
  
  if (deviceId == null) {
    // Gerar novo device_id (mínimo 16 caracteres)
    deviceId = const Uuid().v4();
    await prefs.setString('device_id', deviceId);
  }
  
  return deviceId;
}
```

```typescript
// React Native exemplo
import AsyncStorage from '@react-native-async-storage/async-storage';
import { v4 as uuidv4 } from 'uuid';

export async function getOrCreateDeviceId(): Promise<string> {
  let deviceId = await AsyncStorage.getItem('device_id');
  
  if (!deviceId) {
    deviceId = uuidv4();
    await AsyncStorage.setItem('device_id', deviceId);
  }
  
  return deviceId;
}
```

### 2. Registrar Dispositivo Anônimo

Ao iniciar o app pela primeira vez ou atualizar FCM token:

```dart
// Flutter
Future<void> registerAnonymousDevice() async {
  final deviceId = await getOrCreateDeviceId();
  final fcmToken = await FirebaseMessaging.instance.getToken();
  final position = await Geolocator.getCurrentPosition();
  
  final response = await http.post(
    Uri.parse('$baseUrl/api/v1/devices/register'),
    headers: {'Content-Type': 'application/json'},
    body: jsonEncode({
      'device_id': deviceId,
      'fcm_token': fcmToken,
      'platform': Platform.isIOS ? 'ios' : 'android',
      'model': await DeviceInfo().model,
      'language': 'pt',
      'latitude': position.latitude,
      'longitude': position.longitude,
      'alert_radius_meters': 1000,
    }),
  );
  
  if (response.statusCode == 200) {
    print('Dispositivo registrado com sucesso');
  }
}
```

```typescript
// React Native
import messaging from '@react-native-firebase/messaging';
import Geolocation from '@react-native-community/geolocation';

async function registerAnonymousDevice() {
  const deviceId = await getOrCreateDeviceId();
  const fcmToken = await messaging().getToken();
  
  Geolocation.getCurrentPosition(async (position) => {
    const response = await fetch(`${BASE_URL}/api/v1/devices/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        device_id: deviceId,
        fcm_token: fcmToken,
        platform: Platform.OS,
        language: 'pt',
        latitude: position.coords.latitude,
        longitude: position.coords.longitude,
        alert_radius_meters: 1000,
      }),
    });
    
    const data = await response.json();
    console.log('Dispositivo registrado:', data);
  });
}
```

### 3. Conectar WebSocket Anônimo

```dart
// Flutter
import 'package:web_socket_channel/web_socket_channel.dart';

class AnonymousWebSocketService {
  WebSocketChannel? _channel;
  final String deviceId;
  
  AnonymousWebSocketService(this.deviceId);
  
  void connect() {
    final wsUrl = 'ws://localhost:8000/ws/alerts';
    
    _channel = WebSocketChannel.connect(
      Uri.parse(wsUrl),
    );
    
    // IMPORTANTE: Enviar device_id no primeiro frame
    _channel!.sink.add(jsonEncode({
      'event': 'register',
      'device_id': deviceId,
    }));
    
    _channel!.stream.listen(
      (message) {
        final data = jsonDecode(message);
        _handleNotification(data);
      },
      onError: (error) => print('WebSocket error: $error'),
      onDone: () => print('WebSocket closed'),
    );
  }
  
  void updateLocation(double lat, double lon) {
    _channel?.sink.add(jsonEncode({
      'event': 'update_location',
      'data': {
        'latitude': lat,
        'longitude': lon,
      },
    }));
  }
  
  void _handleNotification(Map<String, dynamic> data) {
    switch (data['event']) {
      case 'new_alert':
        _showAlert(data['data']);
        break;
      case 'report_created':
        _showReport(data['data']);
        break;
    }
  }
  
  void disconnect() {
    _channel?.sink.close();
  }
}
```

### 4. Atualizar Localização Periodicamente

```dart
// Flutter
import 'dart:async';
import 'package:geolocator/geolocator.dart';

class LocationTracker {
  Timer? _timer;
  final AnonymousWebSocketService wsService;
  final http.Client httpClient;
  final String deviceId;
  
  LocationTracker({
    required this.wsService,
    required this.httpClient,
    required this.deviceId,
  });
  
  void startTracking() {
    // Atualizar a cada 30 segundos
    _timer = Timer.periodic(Duration(seconds: 30), (_) async {
      final position = await Geolocator.getCurrentPosition();
      
      // Atualizar via WebSocket (tempo real)
      wsService.updateLocation(
        position.latitude,
        position.longitude,
      );
      
      // Atualizar via HTTP (persistência)
      await httpClient.put(
        Uri.parse('$baseUrl/api/v1/devices/location'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'device_id': deviceId,
          'latitude': position.latitude,
          'longitude': position.longitude,
        }),
      );
    });
  }
  
  void stopTracking() {
    _timer?.cancel();
  }
}
```

---

## API Endpoints para Usuários Anônimos

### 1. Registrar Dispositivo

**Endpoint**: `POST /api/v1/devices/register`

**Headers**: Nenhum (público)

**Request Body**:
```json
{
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "fcm_token": "dQw4w9WgXcQ:APA91b...",
  "platform": "android",
  "model": "Pixel 7",
  "language": "pt",
  "latitude": -8.8383,
  "longitude": 13.2344,
  "alert_radius_meters": 1000
}
```

**Response** (200 OK):
```json
{
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "fcm_token": "dQw4w9WgXcQ:APA91b...",
  "platform": "android",
  "latitude": -8.8383,
  "longitude": 13.2344,
  "alert_radius_meters": 1000,
  "message": "Device registered successfully"
}
```

### 2. Atualizar Localização

**Endpoint**: `PUT /api/v1/devices/location`

**Headers**: Nenhum (público)

**Request Body**:
```json
{
  "device_id": "550e8400-e29b-41d4-a716-446655440000",
  "latitude": -8.8400,
  "longitude": 13.2350
}
```

**Response** (200 OK):
```json
{
  "message": "Location updated successfully"
}
```

### 3. WebSocket Connection

**Endpoint**: `ws://host:port/ws/alerts`

**Headers**:
```
X-Device-ID: 550e8400-e29b-41d4-a716-446655440000
```

**OU** (alternativa):
```
Device-ID: 550e8400-e29b-41d4-a716-446655440000
```

**Mensagens Recebidas**:

1. **Novo Alerta**:
```json
{
  "event": "new_alert",
  "data": {
    "alert_id": "abc-123",
    "message": "🚨 Assalto reportado na área",
    "latitude": -8.8390,
    "longitude": 13.2345,
    "radius": 500
  }
}
```

2. **Novo Report**:
```json
{
  "event": "report_created",
  "data": {
    "report_id": "def-456",
    "message": "📍 Buraco na via",
    "latitude": -8.8395,
    "longitude": 13.2348
  }
}
```

---

## Diferenças entre Usuários Autenticados e Anônimos

| Funcionalidade | Usuário Autenticado | Usuário Anônimo |
|----------------|---------------------|-----------------|
| **Receber Alertas** | ✅ Sim | ✅ Sim |
| **Receber Reports** | ✅ Sim | ✅ Sim |
| **WebSocket** | ✅ JWT Token | ✅ Device ID |
| **Push Notifications** | ✅ Sim | ✅ Sim |
| **Criar Alertas** | ✅ Sim | ❌ Não |
| **Criar Reports** | ✅ Sim | ❌ Não |
| **Histórico** | ✅ Sim | ❌ Não |
| **Perfil** | ✅ Sim | ❌ Não |

---

## Migração de Anônimo para Autenticado

Quando o usuário decide criar uma conta:

```dart
// Flutter
Future<void> migrateToAuthenticatedUser(String email, String password) async {
  final deviceId = await getOrCreateDeviceId();
  final fcmToken = await FirebaseMessaging.instance.getToken();
  
  // 1. Criar conta
  final signupResponse = await http.post(
    Uri.parse('$baseUrl/api/v1/auth/signup'),
    body: jsonEncode({
      'email': email,
      'password': password,
      'name': 'Nome do Usuário',
      // ... outros campos
    }),
  );
  
  // 2. Fazer login
  final loginResponse = await http.post(
    Uri.parse('$baseUrl/api/v1/auth/login'),
    body: jsonEncode({
      'email': email,
      'password': password,
    }),
  );
  
  final jwt = jsonDecode(loginResponse.body)['token'];
  await storage.write(key: 'jwt_token', value: jwt);
  
  // 3. Atualizar FCM token do usuário autenticado
  await http.put(
    Uri.parse('$baseUrl/api/v1/users/me/device'),
    headers: {
      'Authorization': 'Bearer $jwt',
      'Content-Type': 'application/json',
    },
    body: jsonEncode({
      'fcm_token': fcmToken,
      'language': 'pt',
    }),
  );
  
  // 4. Reconectar WebSocket com JWT
  wsService.disconnect();
  wsService.connectAuthenticated(jwt);
}
```

---

## Boas Práticas

### 1. Persistência do Device ID
- ✅ Armazenar em `SharedPreferences` / `AsyncStorage`
- ✅ Nunca regenerar após instalação
- ✅ Mínimo 16 caracteres (UUID recomendado)

### 2. Gerenciamento de Localização
- ✅ Solicitar permissão ao usuário
- ✅ Atualizar a cada 30-60 segundos quando em movimento
- ✅ Parar updates quando app está em background (economia de bateria)

### 3. Conexão WebSocket
- ✅ Implementar reconexão automática
- ✅ Usar exponential backoff em caso de falha
- ✅ Desconectar quando app vai para background

### 4. Push Notifications
- ✅ Atualizar FCM token quando ele mudar
- ✅ Implementar tratamento de notificações no background
- ✅ Sincronizar com servidor após receber push

---

## Exemplo Completo Flutter

```dart
import 'package:flutter/material.dart';
import 'package:geolocator/geolocator.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class AnonymousUserService {
  static const String baseUrl = 'http://localhost:8000';
  late String deviceId;
  late AnonymousWebSocketService wsService;
  
  Future<void> initialize() async {
    // 1. Obter ou criar device ID
    deviceId = await getOrCreateDeviceId();
    
    // 2. Obter permissão de localização
    await _requestLocationPermission();
    
    // 3. Obter FCM token
    final fcmToken = await FirebaseMessaging.instance.getToken();
    
    // 4. Registrar dispositivo no backend
    await registerDevice(fcmToken);
    
    // 5. Conectar WebSocket
    wsService = AnonymousWebSocketService(deviceId);
    wsService.connect();
    
    // 6. Iniciar tracking de localização
    _startLocationTracking();
    
    // 7. Configurar push notifications
    _setupPushNotifications();
  }
  
  Future<String> getOrCreateDeviceId() async {
    final prefs = await SharedPreferences.getInstance();
    String? id = prefs.getString('device_id');
    
    if (id == null) {
      id = const Uuid().v4();
      await prefs.setString('device_id', id);
    }
    
    return id;
  }
  
  Future<void> _requestLocationPermission() async {
    LocationPermission permission = await Geolocator.checkPermission();
    
    if (permission == LocationPermission.denied) {
      permission = await Geolocator.requestPermission();
    }
  }
  
  Future<void> registerDevice(String? fcmToken) async {
    final position = await Geolocator.getCurrentPosition();
    
    final response = await http.post(
      Uri.parse('$baseUrl/api/v1/devices/register'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'device_id': deviceId,
        'fcm_token': fcmToken,
        'platform': Platform.isIOS ? 'ios' : 'android',
        'language': 'pt',
        'latitude': position.latitude,
        'longitude': position.longitude,
        'alert_radius_meters': 1000,
      }),
    );
    
    if (response.statusCode == 200) {
      print('✅ Dispositivo registrado');
    } else {
      print('❌ Erro ao registrar: ${response.body}');
    }
  }
  
  void _startLocationTracking() {
    Timer.periodic(Duration(seconds: 30), (_) async {
      final position = await Geolocator.getCurrentPosition();
      wsService.updateLocation(position.latitude, position.longitude);
    });
  }
  
  void _setupPushNotifications() {
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      print('📩 Notificação recebida: ${message.notification?.title}');
      // Exibir notificação local
    });
  }
}
```

---

## Troubleshooting

### Problema: Não recebo notificações

**Checklist**:
1. ✅ Device ID está correto e persistente?
2. ✅ FCM token está atualizado?
3. ✅ Localização está sendo atualizada?
4. ✅ WebSocket está conectado?
5. ✅ Raio de alerta está configurado (default: 1000m)?

### Problema: WebSocket desconecta frequentemente

**Solução**:
- Implementar reconexão automática
- Verificar rede do dispositivo
- Usar heartbeat/ping a cada 30 segundos

### Problema: Localização não atualiza

**Solução**:
- Verificar permissões de localização
- Confirmar GPS está habilitado
- Verificar chamadas HTTP/WebSocket de update

---

## Segurança

### Limitações de Usuários Anônimos

- ❌ Não podem criar alertas
- ❌ Não podem criar reports
- ❌ Não podem verificar reports
- ❌ Não têm acesso a histórico
- ✅ Apenas recebem notificações passivamente

### Limpeza de Sessões Antigas

Sessões anônimas inativas por mais de **30 dias** são automaticamente removidas.

---

## Suporte

Para dúvidas ou problemas, contate a equipe de desenvolvimento.
