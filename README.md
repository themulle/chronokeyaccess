# Chrono Key Access

Chrono Key Access ist ein nutzerfreundliches, Pin-Code-basiertes Offline-Türschließsystem für für Wiegand Codereader Pin-Eingabe.
Es ermöglicht zeitlich begrenzten Zutritt zu gewähren, ohne eine Online-Verbindung zu benötigen.
Hierzu wird ein sehr sicherer Algorithmus verwendet, welcher die Pin-Codes anhand des Zeitstempels generiert.

## Die Vorteile

- **Bequemer Zugang**: Ermöglicht Zutritt zu definierten Zeiten ohne Online-Anbindung des Schließsystems
- **Sicher und Kontrolliert**: Zutritt nur zu vordefinierten Zeiten möglich.
- **Flexibel und Einfach**: Einheitliche oder separate Pin-Codes für unterschiedliche Tageszeiten.
- **Nutzerfreundlich**: Anwender können die Pin-Codes über Web-Applikationen, APIs oder die Kommandozeile abfragen.
- **Einfache Konfiguration**: Flexible Definition der Zutrittszeiten über weit verbreitete Cron-Syntax.

## Nachteile

- Einmal übermittelter PIN an einen Benutzer kann nicht widerrufen/gesperrt werden.
- Benutzer benötigt eventuell mehrere Pin-Codes (einen pro Tag)
- Manuelle Übertragung von Konfigurationsänderungen auf das Schließsystem erforderlich.

## Hintergrund

Entwickelt, um den Mietern der Kalthalle des SV Ringingen den Zugang nur zu ihren gebuchten Zeiten zu gewähren. Diese Applikation ermöglicht es dem Verwalter der Kalthalle die Zutrittcodes für die gebuchten Zeiten direkt nach der Buchung über WhatsApp oder E-Mail an die Mieter zu senden.

## Beispielkonfigurationen

**Jeden Montag von Januar bis März, von 16:00 bis 20:00 Uhr, denselben Pin-Code verwenden**

```plaintext
CronString: 0 0 16 ? JAN,FEB,MAR MON *
Dauer: 4h
Wechselnder Code: Nein

Mo 02.01.2023 16:00 - 20:00 PIN:6467
Mo 09.01.2023 16:00 - 20:00 PIN:6467
Mo 16.01.2023 16:00 - 20:00 PIN:6467
...
Mo 13.03.2023 16:00 - 20:00 PIN:6467
Mo 20.03.2023 16:00 - 20:00 PIN:6467
Mo 27.03.2023 16:00 - 20:00 PIN:6467
```

**Jeden Montag von Januar bis März, von 16:00 bis 20:00 Uhr, einen wechselnden Pin-Code verwenden**
```plaintext
CronString: 0 0 16 ? JAN,FEB,MAR MON *
Dauer: 4h
Wechselnder Code: Ja

Mo 02.01.2023 16:00 - 20:00 PIN:0976
Mo 09.01.2023 16:00 - 20:00 PIN:4328
Mo 16.01.2023 16:00 - 20:00 PIN:3480
...
Mo 13.03.2023 16:00 - 20:00 PIN:4328
Mo 20.03.2023 16:00 - 20:00 PIN:8432
Mo 27.03.2023 16:00 - 20:00 PIN:5320
```

**Mehr details zur Cron-Syntax:**  
[Syntaxreferenz](https://github.com/gorhill/cronexpr)  
[Online-Editor](https://www.freeformatter.com/cron-expression-generator-quartz.html)  
