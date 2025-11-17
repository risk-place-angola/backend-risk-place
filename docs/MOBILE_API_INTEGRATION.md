# Guia de Integração API - Mobile

Este documento descreve todas as rotas disponíveis na API do Risk Place para integração com aplicações mobile.

## Base URL
```
http://localhost:8080/api/v1
```

## Autenticação

A maioria das rotas requer autenticação via JWT Bearer Token. Inclua o token no header:
```
Authorization: Bearer {seu_token_aqui}
```

---

## 📋 Índice de Rotas

- [Autenticação](#autenticação)
- [Usuário](#usuário)
- [Alertas](#alertas)
- [Relatórios](#relatórios)
- [Riscos](#riscos)
- [WebSocket](#websocket)

---

## 🔐 Autenticação

### 1. Cadastro de Usuário
**Endpoint:** `POST /auth/signup`  
**Autenticação:** Não requerida  

**Request Body:**
```json
{
  "name": "João Silva",
  "email": "joao@example.com",
  "phone": "+244923456789",
  "password": "senha123"
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Possíveis Erros:**
- `400 Bad Request` - Dados inválidos
- `500 Internal Server Error` - Erro no servidor

---

### 2. Login
**Endpoint:** `POST /auth/login`  
**Autenticação:** Não requerida  

**Request Body:**
```json
{
  "email": "joao@example.com",
  "password": "senha123"
}
```

**Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "active_role": "user",
    "name": "João Silva",
    "email": "joao@example.com",
    "role_name": ["user"]
  }
}
```

**Possíveis Erros:**
- `400 Bad Request` - Dados inválidos
- `401 Unauthorized` - Credenciais inválidas
- `403 Forbidden` - Email não verificado

---

### 3. Confirmar Cadastro
**Endpoint:** `POST /auth/confirm`  
**Autenticação:** Não requerida  

**Request Body:**
```json
{
  "email": "joao@example.com",
  "code": "123456"
}
```

**Response (204 No Content)**

**Possíveis Erros:**
- `400 Bad Request` - Código inválido ou expirado
- `404 Not Found` - Usuário não encontrado

---

### 4. Esqueci Minha Senha
**Endpoint:** `POST /auth/password/forgot`  
**Autenticação:** Não requerida  

**Request Body:**
```json
{
  "email": "joao@example.com"
}
```

**Response (200 OK):**
```json
"password reset code sent"
```

**Possíveis Erros:**
- `400 Bad Request` - Email inválido
- `404 Not Found` - Usuário não encontrado

---

### 5. Resetar Senha
**Endpoint:** `POST /auth/password/reset`  
**Autenticação:** Não requerida  

**Request Body:**
```json
{
  "email": "joao@example.com",
  "password": "novaSenha123"
}
```

**Response (200 OK):**
```json
"password reset successfully"
```

**Possíveis Erros:**
- `400 Bad Request` - Código inválido
- `404 Not Found` - Usuário não encontrado

---

## 👤 Usuário

### 6. Obter Perfil do Usuário
**Endpoint:** `GET /users/me`  
**Autenticação:** Requerida (JWT)  

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "João Silva",
  "email": "joao@example.com",
  "phone": "+244923456789",
  "nif": "123456789",
  "role_name": ["user"],
  "address": {
    "Country": "Angola",
    "Province": "Luanda",
    "Municipality": "Luanda",
    "Neighborhood": "Talatona",
    "ZipCode": "12345"
  }
}
```

**Possíveis Erros:**
- `401 Unauthorized` - Token inválido ou ausente
- `404 Not Found` - Usuário não encontrado

---

## 🚨 Alertas

### 7. Criar Alerta
**Endpoint:** `POST /alerts`  
**Autenticação:** Requerida (JWT)  

**Request Body:**
```json
{
  "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
  "risk_topic_id": "550e8400-e29b-41d4-a716-446655440002",
  "message": "Assalto em andamento na área",
  "latitude": -8.8383,
  "longitude": 13.2344,
  "radius": 500.0,
  "severity": "high"
}
```

**Campos:**
- `risk_type_id` (string, obrigatório): UUID do tipo de risco
- `risk_topic_id` (string, obrigatório): UUID do tópico de risco
- `message` (string, obrigatório): Mensagem do alerta
- `latitude` (number, obrigatório): Latitude da localização
- `longitude` (number, obrigatório): Longitude da localização
- `radius` (number, obrigatório): Raio de alcance em metros
- `severity` (string, obrigatório): Gravidade (low, medium, high, critical)

**Response (201 Created):**
```json
{
  "status": "alert triggered"
}
```

**Possíveis Erros:**
- `400 Bad Request` - Dados inválidos
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao processar

**Nota:** O alerta é enviado via WebSocket para todos os usuários conectados no raio especificado.

---

## 📍 Relatórios

### 8. Criar Relatório
**Endpoint:** `POST /reports`  
**Autenticação:** Requerida (JWT)  

**Request Body:**
```json
{
  "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
  "risk_topic_id": "550e8400-e29b-41d4-a716-446655440002",
  "description": "Buraco grande na via principal",
  "latitude": -8.8383,
  "longitude": 13.2344,
  "province": "Luanda",
  "municipality": "Luanda",
  "neighborhood": "Talatona",
  "address": "Rua Principal, próximo ao Shopping",
  "image_url": "https://example.com/image.jpg"
}
```

**Campos:**
- `risk_type_id` (string, obrigatório): UUID do tipo de risco
- `risk_topic_id` (string, obrigatório): UUID do tópico de risco
- `description` (string, obrigatório): Descrição do problema
- `latitude` (number, obrigatório): Latitude da localização
- `longitude` (number, obrigatório): Longitude da localização
- `province` (string, opcional): Província
- `municipality` (string, opcional): Município
- `neighborhood` (string, opcional): Bairro
- `address` (string, opcional): Endereço completo
- `image_url` (string, opcional): URL da imagem do problema

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
  "risk_topic_id": "550e8400-e29b-41d4-a716-446655440002",
  "description": "Buraco grande na via principal",
  "latitude": -8.8383,
  "longitude": 13.2344,
  "province": "Luanda",
  "municipality": "Luanda",
  "neighborhood": "Talatona",
  "address": "Rua Principal, próximo ao Shopping",
  "image_url": "https://example.com/image.jpg",
  "status": "pending",
  "created_at": "2025-11-17T10:30:00Z",
  "updated_at": "2025-11-17T10:30:00Z"
}
```

**Possíveis Erros:**
- `400 Bad Request` - Dados inválidos
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao criar relatório

---

### 9. Listar Relatórios Próximos
**Endpoint:** `GET /reports/nearby`  
**Autenticação:** Requerida (JWT)  

**Query Parameters:**
- `lat` (obrigatório): Latitude do ponto de referência
- `lon` (obrigatório): Longitude do ponto de referência
- `radius` (opcional): Raio em metros (padrão: 500)

**Exemplo:**
```
GET /reports/nearby?lat=-8.8383&lon=13.2344&radius=1000
```

**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
    "risk_topic_id": "550e8400-e29b-41d4-a716-446655440002",
    "description": "Buraco grande na via principal",
    "latitude": -8.8383,
    "longitude": 13.2344,
    "province": "Luanda",
    "municipality": "Luanda",
    "neighborhood": "Talatona",
    "address": "Rua Principal, próximo ao Shopping",
    "image_url": "https://example.com/image.jpg",
    "status": "pending",
    "created_at": "2025-11-17T10:30:00Z",
    "updated_at": "2025-11-17T10:30:00Z"
  }
]
```

**Possíveis Status:**
- `pending` - Pendente de verificação
- `verified` - Verificado
- `resolved` - Resolvido

**Possíveis Erros:**
- `400 Bad Request` - Parâmetros inválidos
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao buscar relatórios

---

### 10. Verificar Relatório
**Endpoint:** `POST /reports/{id}/verify`  
**Autenticação:** Requerida (JWT)  

**URL Parameters:**
- `id`: UUID do relatório

**Request Body:**
```json
{
  "moderator_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (200 OK):**
```json
{
  "status": "verified",
  "report_id": "550e8400-e29b-41d4-a716-446655440003"
}
```

**Possíveis Erros:**
- `400 Bad Request` - ID inválido
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao verificar

---

### 11. Resolver Relatório
**Endpoint:** `POST /reports/{id}/resolve`  
**Autenticação:** Requerida (JWT)  

**URL Parameters:**
- `id`: UUID do relatório

**Request Body:**
```json
{
  "moderator_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (200 OK):**
```json
{
  "status": "resolved",
  "report_id": "550e8400-e29b-41d4-a716-446655440003"
}
```

**Possíveis Erros:**
- `400 Bad Request` - ID inválido
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao resolver

---

## ⚠️ Riscos

### 12. Listar Tipos de Risco
**Endpoint:** `GET /risks/types`  
**Autenticação:** Requerida (JWT)  

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Crime",
      "description": "Atividades criminosas e segurança pública",
      "default_radius": 500,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "Infraestrutura",
      "description": "Problemas relacionados à infraestrutura urbana",
      "default_radius": 1000,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

**Possíveis Erros:**
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao buscar tipos

---

### 13. Listar Tópicos de Risco
**Endpoint:** `GET /risks/topics`  
**Autenticação:** Requerida (JWT)  

**Query Parameters (opcionais):**
- `risk_type_id`: UUID para filtrar tópicos por tipo de risco

**Exemplo sem filtro:**
```
GET /risks/topics
```

**Exemplo com filtro:**
```
GET /risks/topics?risk_type_id=550e8400-e29b-41d4-a716-446655440001
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440010",
      "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Assalto",
      "description": "Roubo à mão armada",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440011",
      "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Furto",
      "description": "Furto sem violência",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440012",
      "risk_type_id": "550e8400-e29b-41d4-a716-446655440002",
      "name": "Buraco na rua",
      "description": "Buracos e problemas no asfalto",
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z"
    }
  ]
}
```

**Possíveis Erros:**
- `401 Unauthorized` - Token inválido
- `500 Internal Server Error` - Erro ao buscar tópicos

---

## 🔌 WebSocket

### 14. Conexão WebSocket para Alertas
**Endpoint:** `ws://localhost:8080/ws/alerts`  
**Protocolo:** WebSocket  

**Autenticação:** Token JWT deve ser enviado após a conexão

**Fluxo de Conexão:**

1. **Conectar ao WebSocket:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/alerts');
```

2. **Enviar token de autenticação após conexão:**
```json
{
  "type": "auth",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

3. **Enviar localização para receber alertas próximos:**
```json
{
  "type": "location",
  "latitude": -8.8383,
  "longitude": 13.2344
}
```

**Mensagens Recebidas:**

**Alerta de Proximidade:**
```json
{
  "type": "alert",
  "data": {
    "risk_type_id": "550e8400-e29b-41d4-a716-446655440001",
    "risk_topic_id": "550e8400-e29b-41d4-a716-446655440002",
    "message": "Assalto em andamento na área",
    "latitude": -8.8383,
    "longitude": 13.2344,
    "radius": 500.0,
    "severity": "high",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "timestamp": "2025-11-17T10:30:00Z"
  }
}
```

**Erro de Autenticação:**
```json
{
  "type": "error",
  "message": "authentication failed"
}
```

**Nota:** Mantenha a conexão WebSocket ativa e atualize a localização periodicamente para receber alertas em tempo real.

Para mais detalhes sobre a implementação do WebSocket, consulte: [MOBILE_WEBSOCKET_INTEGRATION.md](./MOBILE_WEBSOCKET_INTEGRATION.md)

---

## 🔄 Health Check

### 15. Verificar Status da API
**Endpoint:** `GET /health`  
**Autenticação:** Não requerida  

**Response (200 OK):**
```
OK
```

---

## 📝 Códigos de Status HTTP

- `200 OK` - Requisição bem-sucedida
- `201 Created` - Recurso criado com sucesso
- `204 No Content` - Requisição bem-sucedida sem conteúdo de retorno
- `400 Bad Request` - Dados inválidos ou mal formatados
- `401 Unauthorized` - Autenticação necessária ou falhou
- `403 Forbidden` - Sem permissão para acessar o recurso
- `404 Not Found` - Recurso não encontrado
- `500 Internal Server Error` - Erro interno do servidor

---

## 📱 Fluxo Recomendado para Mobile

### Primeiro Acesso:
1. Cadastro (`POST /auth/signup`)
2. Confirmar email (`POST /auth/confirm`)
3. Login (`POST /auth/login`)
4. Armazenar tokens
5. Conectar ao WebSocket (`/ws/alerts`)

### Uso Regular:
1. Verificar token armazenado
2. Se válido, conectar ao WebSocket
3. Enviar localização atual
4. Buscar tipos e tópicos de risco (`GET /risks/types`, `GET /risks/topics`)
5. Listar relatórios próximos (`GET /reports/nearby`)
6. Criar alertas/relatórios quando necessário

### Gestão de Tokens:
- Armazene o `access_token` de forma segura
- Renove usando o `refresh_token` quando expirar
- Implemente logout limpando tokens armazenados

---

## 🛠️ Exemplos de Implementação

### Flutter / Dart
```dart
// Login
final response = await http.post(
  Uri.parse('http://localhost:8080/api/v1/auth/login'),
  headers: {'Content-Type': 'application/json'},
  body: jsonEncode({
    'email': 'joao@example.com',
    'password': 'senha123',
  }),
);

if (response.statusCode == 200) {
  final data = jsonDecode(response.body);
  final token = data['access_token'];
  // Armazenar token
}
```

### React Native / JavaScript
```javascript
// Criar Alerta
const createAlert = async (alertData, token) => {
  const response = await fetch('http://localhost:8080/api/v1/alerts', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(alertData)
  });
  
  return await response.json();
};
```

### Swift / iOS
```swift
// Buscar Relatórios Próximos
func fetchNearbyReports(lat: Double, lon: Double, token: String) {
    let url = URL(string: "http://localhost:8080/api/v1/reports/nearby?lat=\(lat)&lon=\(lon)")!
    var request = URLRequest(url: url)
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
    
    URLSession.shared.dataTask(with: request) { data, response, error in
        // Processar resposta
    }.resume()
}
```

---

## 🔍 Documentação Adicional

- **Swagger UI:** [http://localhost:8080/docs/](http://localhost:8080/docs/)
- **WebSocket Guide:** [MOBILE_WEBSOCKET_INTEGRATION.md](./MOBILE_WEBSOCKET_INTEGRATION.md)
- **Notification Guide:** [WEBSOCKET_NOTIFICATION_GUIDE.md](./WEBSOCKET_NOTIFICATION_GUIDE.md)

---

## 📞 Suporte

Para dúvidas ou problemas na integração, consulte a documentação completa ou entre em contato com a equipe de desenvolvimento.

**Última Atualização:** 17 de Novembro de 2025
