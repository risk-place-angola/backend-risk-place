# Risk Place Angola - Documentation Index

Bem-vindo à documentação técnica do backend Risk Place Angola.

## 📚 Guias de Integração

### Mobile Development

#### WebSocket & Real-Time
- **[WEBSOCKET_GUIDE.md](./WEBSOCKET_GUIDE.md)** - Guia completo de integração WebSocket
  - Configuração de conexão
  - Protocolo de mensagens
  - Tipos de eventos
  - Autenticação (JWT e anônimos)
  - Tratamento de erros e reconexão
  - Best practices

#### REST API
- **[MOBILE_API_INTEGRATION.md](./MOBILE_API_INTEGRATION.md)** - Documentação completa da API REST
  - Endpoints de autenticação
  - Gerenciamento de usuários
  - Alertas e reports
  - Tipos e tópicos de riscos
  - Exemplos de requisições e respostas

#### Framework Específico
- **[FLUTTER_INTEGRATION_GUIDE.md](./FLUTTER_INTEGRATION_GUIDE.md)** - Implementação Flutter passo a passo
  - Setup de dependências
  - Modelos de dados (DTOs)
  - Services (API, WebSocket, Location)
  - Gerenciador de usuários anônimos
  - Exemplos completos de código
  - Troubleshooting

### Usuários Anônimos

- **[ANONYMOUS_USER_GUIDE.md](./ANONYMOUS_USER_GUIDE.md)** - Documentação técnica completa
  - Arquitetura do sistema
  - Implementação backend
  - Integração mobile
  - Fluxos de notificação
  - Casos de uso

## 🏗️ Arquitetura

### Diagramas

- **[diagram/ANONYMOUS_USERS_ARCHITECTURE.md](./diagram/ANONYMOUS_USERS_ARCHITECTURE.md)** - Arquitetura de usuários anônimos
- **HighLevelArchitecture.svg** - Arquitetura geral do sistema
- **CleanArchitectureLayers.svg** - Camadas da Clean Architecture
- **TriggerAlertFlow.svg** - Fluxo de criação de alertas
- **ReportFlow.svg** - Fluxo de reports
- **UserAuthenticationFlow.svg** - Fluxo de autenticação
- **NotificationFlow.svg** - Fluxo de notificações

## 🎯 Públicos-Alvo

### Para Desenvolvedores Mobile
1. Comece com **[MOBILE_API_INTEGRATION.md](./MOBILE_API_INTEGRATION.md)** para entender os endpoints
2. Leia **[WEBSOCKET_GUIDE.md](./WEBSOCKET_GUIDE.md)** para implementar notificações em tempo real
3. Para Flutter, use **[FLUTTER_INTEGRATION_GUIDE.md](./FLUTTER_INTEGRATION_GUIDE.md)** como referência de implementação
4. Se implementar usuários anônimos, consulte **[ANONYMOUS_USER_GUIDE.md](./ANONYMOUS_USER_GUIDE.md)**

### Para Desenvolvedores Backend
1. Leia **[WEBSOCKET_GUIDE.md](./WEBSOCKET_GUIDE.md)** para entender a arquitetura WebSocket
2. Consulte **[ANONYMOUS_USER_GUIDE.md](./ANONYMOUS_USER_GUIDE.md)** para detalhes do sistema de anônimos
3. Revise os diagramas de arquitetura em **diagram/**

### Para QA/Testes
1. **[WEBSOCKET_GUIDE.md](./WEBSOCKET_GUIDE.md)** - Seção de testes com websocat
2. **[MOBILE_API_INTEGRATION.md](./MOBILE_API_INTEGRATION.md)** - Exemplos de curl para testes

## 🔄 Quick Links

| Preciso de... | Documento |
|---------------|-----------|
| Conectar WebSocket | [WEBSOCKET_GUIDE.md](./WEBSOCKET_GUIDE.md) |
| Listar endpoints da API | [MOBILE_API_INTEGRATION.md](./MOBILE_API_INTEGRATION.md) |
| Implementar em Flutter | [FLUTTER_INTEGRATION_GUIDE.md](./FLUTTER_INTEGRATION_GUIDE.md) |
| Suportar usuários anônimos | [ANONYMOUS_USER_GUIDE.md](./ANONYMOUS_USER_GUIDE.md) |
| Entender arquitetura | [diagram/](./diagram/) |

## 📝 Changelog

### v1.0.0 (Novembro 2025)
- ✅ Sistema de usuários anônimos
- ✅ WebSocket com suporte a device_id
- ✅ Push notifications via FCM
- ✅ Geolocation com Redis
- ✅ Documentação completa consolidada

## 🤝 Contribuindo

Encontrou algo faltando ou incorreto na documentação?
1. Abra uma issue no repositório
2. Ou envie um PR com correções/melhorias

## 📞 Suporte

- **Discord**: [Join our server](https://discord.gg/s2Nk4xYV)
- **GitHub Issues**: [Report issues](https://github.com/risk-place-angola/backend-risk-place/issues)

---

**Última atualização**: Novembro 17, 2025
