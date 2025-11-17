# WebSocket Real-Time Notification System
## Risk Place Angola

**Version**: 1.0.0  
**Last Updated**: November 17, 2025  
**Target**: Mobile & Web Applications

---

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Connection Setup](#connection-setup)
- [Message Protocol](#message-protocol)
- [Event Types](#event-types)
- [Location Updates](#location-updates)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [Testing](#testing)

---

## Overview

O backend Risk Place Angola implementa um sistema de notificações em tempo real via WebSockets que permite:

- ✅ Receber alertas instantâneos sobre riscos próximos
- ✅ Receber reports de incidentes na vizinhança
- ✅ Atualizar localização do usuário em tempo real
- ✅ Manter conexão persistente com reconexão automática
- ✅ Suporte para usuários autenticados e anônimos

### Fluxo de Funcionamento

```
┌─────────────┐
│ Mobile App  │
└──────┬──────┘
       │ 1. Autenticar (JWT) ou Registrar (device_id)
       ▼
┌─────────────┐
│ Backend API │
└──────┬──────┘
       │ 2. Conectar WebSocket
       │    ws://host/ws/alerts
       │    Header: Authorization: Bearer <JWT>
       │    OU Header: X-Device-ID: <device_id>
       ▼
┌─────────────┐
│ WebSocket   │
│    Hub      │
└──────┬──────┘
       │ 3. Registrar Cliente
       │ 4. Enviar/Receber Mensagens
       ▼
┌─────────────┐
│   Active    │
│  Session    │
└─────────────┘
```

---

## Architecture

### Componentes Principais

#### 1. WebSocket Hub
- **Localização**: `internal/adapter/websocket/websocket_hub.go`
- **Função**: Gerencia todas as conexões ativas
- **Responsabilidades**:
  - Registrar/desregistrar clientes
  - Broadcast de mensagens
  - Processar atualizações de localização

#### 2. WebSocket Client
- **Localização**: `internal/adapter/websocket/websocket_client.go`
- **Função**: Representa uma conexão individual
- **Responsabilidades**:
  - Gerenciar canal de mensagens
  - Heartbeat/keep-alive
  - Tratamento de erros

#### 3. Location Store (Redis)
- **Localização**: `internal/infra/location/redis_location_store.go`
- **Função**: Armazenamento geoespacial
- **Responsabilidades**:
  - Indexar localizações com Redis GEOADD
  - Buscar usuários em raio com GEOSEARCH
  - Key: `user_locations`

#### 4. Event Dispatcher
- **Localização**: `internal/domain/event/dispatcher.go`
- **Função**: Coordenar eventos do sistema
- **Responsabilidades**:
  - Disparar notificações
  - Integrar com FCM para offline users

### Fluxo de Notificações

```
┌──────────────┐
│ Alert/Report │
│   Created    │
└──────┬───────┘
       │ 1. Domain Event
       ▼
┌──────────────┐
│   Event      │
│  Dispatcher  │
└──────┬───────┘
       │ 2. Query Redis
       ▼
┌──────────────┐         ┌──────────────┐
│ Redis Geo    │         │ WebSocket    │
│ Find Nearby  │────────>│   Broadcast  │
└──────────────┘    3.   └──────┬───────┘
                                 │ 4. Send to Clients
                                 ▼
                         ┌──────────────┐
                         │   Mobile     │
                         │    Apps      │
                         └──────────────┘
```

---

## Connection Setup

### Endpoint Information

| Environment | WebSocket URL | Protocol |
|-------------|---------------|----------|
| Development | `ws://localhost:8000/ws/alerts` | ws:// |
| Production  | `wss://api.riskplace.com/ws/alerts` | wss:// (TLS) |

### Connection Methods

#### Opção 1: Usuário Autenticado (JWT)

```http
GET ws://localhost:8000/ws/alerts HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Upgrade: websocket
Connection: Upgrade
```

#### Opção 2: Usuário Anônimo (Device ID)

```http
GET ws://localhost:8000/ws/alerts HTTP/1.1
X-Device-ID: 550e8400-e29b-41d4-a716-446655440000
Upgrade: websocket
Connection: Upgrade
```

### Exemplo Flutter

```dart
import 'package:web_socket_channel/web_socket_channel.dart';

// Com JWT
final channel = WebSocketChannel.connect(
  Uri.parse('ws://localhost:8000/ws/alerts'),
);

// Adicionar header após conexão (limitação do package)
// Solução: Enviar device_id na primeira mensagem

final deviceId = '550e8400-e29b-41d4-a716-446655440000';
channel.sink.add(jsonEncode({
  'event': 'register',
  'device_id': deviceId,
}));

// Escutar mensagens
channel.stream.listen((message) {
  final data = jsonDecode(message);
  print('Received: ${data['event']}');
});
```

---

## Message Protocol

Todas as mensagens seguem o formato JSON:

```json
{
  "event": "event_type",
  "data": { }
}
```

### Mensagens do Cliente → Servidor

#### 1. Registrar Dispositivo (Anônimo)

```json
{
  "event": "register",
  "device_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 2. Atualizar Localização

```json
{
  "event": "update_location",
  "data": {
    "latitude": -8.8390,
    "longitude": 13.2345
  }
}
```

**Resposta**:
```json
{
  "event": "location_updated",
  "data": {
    "status": "ok"
  }
}
```

#### 3. Heartbeat (Keep-Alive)

```json
{
  "event": "ping"
}
```

**Resposta**:
```json
{
  "event": "pong"
}
```

### Mensagens do Servidor → Cliente

#### 1. Novo Alerta

```json
{
  "event": "new_alert",
  "data": {
    "alert_id": "660f9511-f3ac-52e5-b827-557766551111",
    "message": "🚨 Assalto reportado na área - Zona de Maianga",
    "latitude": -8.8390,
    "longitude": 13.2345,
    "radius": 500.0,
    "severity": "high",
    "created_at": "2025-11-17T14:30:00Z"
  }
}
```

#### 2. Novo Report

```json
{
  "event": "report_created",
  "data": {
    "report_id": "770fa622-g4bd-63f6-c938-668877662222",
    "message": "📍 Buraco grande na via - Avenida 4 de Fevereiro",
    "latitude": -8.8395,
    "longitude": 13.2348,
    "risk_type": "infrastructure",
    "created_at": "2025-11-17T14:35:00Z"
  }
}
```

#### 3. Erro

```json
{
  "event": "error",
  "data": {
    "code": "UNAUTHORIZED",
    "message": "Invalid authentication token"
  }
}
```

---

## Event Types

### Eventos Recebidos pelo Cliente

| Event | Descrição | Quando Ocorre |
|-------|-----------|---------------|
| `new_alert` | Alerta de risco criado | Quando ERCE/ERFCE cria alerta próximo |
| `report_created` | Report de incidente | Quando usuário reporta problema próximo |
| `location_updated` | Confirmação de localização | Após `update_location` bem-sucedido |
| `pong` | Resposta ao heartbeat | Após cliente enviar `ping` |
| `error` | Erro na operação | Quando ocorre falha |

### Eventos Enviados pelo Cliente

| Event | Descrição | Frequência Recomendada |
|-------|-----------|------------------------|
| `register` | Registrar device_id | Uma vez ao conectar (anônimos) |
| `update_location` | Atualizar posição GPS | A cada 30-60 segundos ou mudança significativa |
| `ping` | Keep-alive | A cada 30 segundos |

---

## Location Updates

### Quando Atualizar

✅ **Recomendado**:
- A cada 30-60 segundos se app estiver em foreground
- Quando usuário se mover > 50 metros
- Ao abrir o app

❌ **Evitar**:
- Atualizações a cada < 10 segundos (sobrecarga)
- Atualizações com app em background (bateria)
- Atualizações sem mudança de localização

### Exemplo de Implementação

```dart
import 'package:geolocator/geolocator.dart';

Timer.periodic(Duration(seconds: 30), (_) async {
  final position = await Geolocator.getCurrentPosition();
  
  // Via WebSocket
  channel.sink.add(jsonEncode({
    'event': 'update_location',
    'data': {
      'latitude': position.latitude,
      'longitude': position.longitude,
    },
  }));
  
  // Também via HTTP (persistência)
  await apiService.updateDeviceLocation(
    deviceId: deviceId,
    latitude: position.latitude,
    longitude: position.longitude,
  );
});
```

---

## Authentication

### Usuários Autenticados

1. **Login via API**:
   ```bash
   POST /api/v1/auth/login
   Content-Type: application/json
   
   {
     "email": "user@example.com",
     "password": "senha123"
   }
   ```

2. **Receber JWT**:
   ```json
   {
     "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
     "user": { "id": "...", "name": "..." }
   }
   ```

3. **Conectar WebSocket**:
   - Header: `Authorization: Bearer <token>`

### Usuários Anônimos

1. **Gerar Device ID**:
   ```dart
   import 'package:uuid/uuid.dart';
   final deviceId = Uuid().v4();
   ```

2. **Registrar Dispositivo**:
   ```bash
   POST /api/v1/devices/register
   Content-Type: application/json
   
   {
     "device_id": "550e8400-e29b-41d4-a716-446655440000",
     "fcm_token": "...",
     "platform": "android",
     "latitude": -8.8383,
     "longitude": 13.2344
   }
   ```

3. **Conectar WebSocket**:
   - Header: `X-Device-ID: <device_id>`
   - OU enviar mensagem `register` após conectar

---

## Error Handling

### Códigos de Erro

| Code | Descrição | Ação |
|------|-----------|------|
| `UNAUTHORIZED` | Token inválido ou expirado | Reautenticar |
| `INVALID_MESSAGE` | Formato JSON inválido | Corrigir mensagem |
| `LOCATION_REQUIRED` | Localização não fornecida | Enviar GPS |
| `CONNECTION_LIMIT` | Muitas conexões simultâneas | Aguardar e reconectar |

### Reconexão Automática

```dart
class WebSocketService {
  int _reconnectAttempts = 0;
  Timer? _reconnectTimer;
  
  void _scheduleReconnect() {
    _reconnectAttempts++;
    
    // Exponential backoff: 2s, 4s, 8s, 16s, max 60s
    final delay = Duration(
      seconds: (2 * _reconnectAttempts).clamp(2, 60),
    );
    
    _reconnectTimer = Timer(delay, () {
      if (!_isConnected) {
        connect();
      }
    });
  }
  
  void _handleError(error) {
    print('WebSocket error: $error');
    _isConnected = false;
    _scheduleReconnect();
  }
}
```

---

## Best Practices

### ✅ Faça

1. **Implementar reconexão automática**
   - Usar exponential backoff
   - Limitar tentativas (max 10)

2. **Enviar heartbeat regularmente**
   - Intervalo: 30 segundos
   - Detectar conexões mortas

3. **Otimizar atualizações de localização**
   - Apenas mudanças significativas (>50m)
   - Respeitar intervalo mínimo (30s)

4. **Tratar erros gracefully**
   - Exibir mensagens amigáveis
   - Log para debug

5. **Persistir device_id**
   - SharedPreferences (Flutter)
   - AsyncStorage (React Native)

### ❌ Evite

1. **Não spam de mensagens**
   - Evitar envios < 10 segundos

2. **Não ignorar erros**
   - Sempre tratar eventos `error`

3. **Não manter múltiplas conexões**
   - Uma conexão por app

4. **Não enviar dados sensíveis**
   - Apenas informações necessárias

---

## Testing

### Teste Manual com `websocat`

```bash
# Instalar websocat
brew install websocat  # macOS
# ou
cargo install websocat  # Rust

# Conectar
websocat ws://localhost:8000/ws/alerts \
  -H "X-Device-ID: 550e8400-e29b-41d4-a716-446655440000"

# Enviar mensagens
{"event":"update_location","data":{"latitude":-8.8390,"longitude":13.2345}}
{"event":"ping"}
```

### Teste de Notificações

1. **Criar alerta próximo**:
   ```bash
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

2. **Verificar recebimento**:
   ```
   Received: {"event":"new_alert","data":{...}}
   ```

### Teste de Reconexão

```dart
// Simular desconexão
channel.sink.close();

// Aguardar reconexão automática
await Future.delayed(Duration(seconds: 5));

// Verificar se reconectou
expect(websocketService.isConnected, true);
```

---

## Troubleshooting

### Problema: WebSocket não conecta

**Soluções**:
- ✅ Verificar URL (ws:// vs wss://)
- ✅ Android emulator: usar `10.0.2.2` em vez de `localhost`
- ✅ iOS simulator: usar `localhost` funciona
- ✅ Verificar firewall/proxy

### Problema: Não recebe notificações

**Soluções**:
- ✅ Confirmar que localização foi atualizada
- ✅ Verificar raio do alerta (deve cobrir sua posição)
- ✅ Checar logs do backend
- ✅ Verificar se conexão está ativa

### Problema: Desconexões frequentes

**Soluções**:
- ✅ Implementar heartbeat
- ✅ Verificar qualidade da rede
- ✅ Aumentar intervalo de reconexão
- ✅ Verificar logs de erro

---

## References

- [Backend Repository](https://github.com/risk-place-angola/backend-risk-place)
- [ANONYMOUS_USER_GUIDE.md](./ANONYMOUS_USER_GUIDE.md)
- [FLUTTER_INTEGRATION_GUIDE.md](./FLUTTER_INTEGRATION_GUIDE.md)
- [MOBILE_API_INTEGRATION.md](./MOBILE_API_INTEGRATION.md)

---

**Versão**: 1.0.0  
**Última Atualização**: Novembro 17, 2025  
**Contato**: Backend Team
