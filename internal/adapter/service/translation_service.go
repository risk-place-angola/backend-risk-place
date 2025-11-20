package service

type Language string

const (
	LanguagePortuguese Language = "pt"
	LanguageEnglish    Language = "en"
)

type NotificationMessage struct {
	Title string
	Body  string
}

type TranslationService struct {
	messages map[string]map[Language]NotificationMessage
}

func NewTranslationService() *TranslationService {
	ts := &TranslationService{
		messages: make(map[string]map[Language]NotificationMessage),
	}
	ts.loadMessages()
	return ts
}

func (ts *TranslationService) loadMessages() {
	ts.messages["alert_created"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🚨 Alerta de Segurança",
			Body:  "Um novo alerta foi criado na sua área",
		},
		LanguageEnglish: {
			Title: "🚨 Safety Alert",
			Body:  "A new alert has been created in your area",
		},
	}

	ts.messages["alert_created_crime"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🚨 Alerta de Crime",
			Body:  "Actividade criminosa reportada próxima de si",
		},
		LanguageEnglish: {
			Title: "🚨 Crime Alert",
			Body:  "Criminal activity reported near you",
		},
	}

	ts.messages["alert_created_accident"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🚨 Alerta de Acidente",
			Body:  "Acidente reportado na sua área",
		},
		LanguageEnglish: {
			Title: "🚨 Accident Alert",
			Body:  "Accident reported in your area",
		},
	}

	ts.messages["alert_created_fire"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🔥 Alerta de Incêndio",
			Body:  "Incêndio reportado próximo de si",
		},
		LanguageEnglish: {
			Title: "🔥 Fire Alert",
			Body:  "Fire reported near you",
		},
	}

	ts.messages["alert_created_natural_disaster"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "⚠️ Alerta de Desastre Natural",
			Body:  "Desastre natural na sua região",
		},
		LanguageEnglish: {
			Title: "⚠️ Natural Disaster Alert",
			Body:  "Natural disaster in your region",
		},
	}

	ts.messages["alert_created_violence"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "⚠️ Alerta de Violência",
			Body:  "Incidente violento reportado na área",
		},
		LanguageEnglish: {
			Title: "⚠️ Violence Alert",
			Body:  "Violent incident reported in the area",
		},
	}

	ts.messages["alert_created_health"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🏥 Alerta de Emergência Médica",
			Body:  "Emergência médica na sua área",
		},
		LanguageEnglish: {
			Title: "🏥 Medical Emergency Alert",
			Body:  "Medical emergency in your area",
		},
	}

	ts.messages["alert_created_infrastructure"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "⚠️ Alerta de Infraestrutura",
			Body:  "Problema de infraestrutura na área",
		},
		LanguageEnglish: {
			Title: "⚠️ Infrastructure Alert",
			Body:  "Infrastructure issue in the area",
		},
	}

	ts.messages["alert_created_environment"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🌍 Alerta Ambiental",
			Body:  "Risco ambiental detectado",
		},
		LanguageEnglish: {
			Title: "🌍 Environmental Alert",
			Body:  "Environmental hazard detected",
		},
	}

	ts.messages["alert_created_public_safety"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "👮 Alerta de Segurança Pública",
			Body:  "Situação de segurança na sua área",
		},
		LanguageEnglish: {
			Title: "👮 Public Safety Alert",
			Body:  "Safety situation in your area",
		},
	}

	ts.messages["alert_created_traffic"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🚦 Alerta de Trânsito",
			Body:  "Problema de trânsito na área",
		},
		LanguageEnglish: {
			Title: "🚦 Traffic Alert",
			Body:  "Traffic issue in the area",
		},
	}

	ts.messages["alert_created_urban_issue"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "🏙️ Alerta Urbano",
			Body:  "Problema urbano reportado",
		},
		LanguageEnglish: {
			Title: "🏙️ Urban Alert",
			Body:  "Urban issue reported",
		},
	}

	ts.messages["report_created"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "📍 Novo Relato",
			Body:  "Novo relato de risco na sua área",
		},
		LanguageEnglish: {
			Title: "📍 New Report",
			Body:  "New risk report in your area",
		},
	}

	ts.messages["report_verified"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "✅ Relato Verificado",
			Body:  "Seu relato foi verificado pelas autoridades",
		},
		LanguageEnglish: {
			Title: "✅ Report Verified",
			Body:  "Your report has been verified by authorities",
		},
	}

	ts.messages["report_resolved"] = map[Language]NotificationMessage{
		LanguagePortuguese: {
			Title: "✅ Relato Resolvido",
			Body:  "O relato na sua área foi resolvido",
		},
		LanguageEnglish: {
			Title: "✅ Report Resolved",
			Body:  "The report in your area has been resolved",
		},
	}
}

func (ts *TranslationService) GetMessage(key string, lang Language, riskType string) NotificationMessage {
	if riskType != "" {
		specificKey := key + "_" + riskType
		if msgs, exists := ts.messages[specificKey]; exists {
			if msg, ok := msgs[lang]; ok {
				return msg
			}
		}
	}

	if msgs, exists := ts.messages[key]; exists {
		if msg, ok := msgs[lang]; ok {
			return msg
		}
		if msg, ok := msgs[LanguagePortuguese]; ok {
			return msg
		}
	}

	return NotificationMessage{
		Title: "Risk Place",
		Body:  "New notification",
	}
}

func (ts *TranslationService) ParseLanguage(lang string) Language {
	switch lang {
	case "en", "EN", "english", "English":
		return LanguageEnglish
	case "pt", "PT", "portuguese", "Portuguese":
		return LanguagePortuguese
	default:
		return LanguagePortuguese
	}
}
