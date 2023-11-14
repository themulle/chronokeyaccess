# Chrono Key Access

Pin-Code basiertes offline Türschließsystem für Codereader.
Das System ermöglicht es zeitlich beschränkten Zutritt zu gewähren ohne dass das Schließsystem online erreichbar ist.
Hierzu werden die PinCodes anhand des aktuellen Zeitstempels errechnet.
Per Konfigurationsdatei werden Zutrittszeiten definiert (z.B. jeden Monatg zwischen 16:00 und 22:00).
Das Programm berechnet hierzu anhand eines Geheimschlüssel den Code der zum entsprechenden Zeitpunkt gültig ist.

**Vorteile:**
- Das Schließsystem muss nicht online sein
- Zugriff für Dritte nur zu den ihnen erlaubten Zeiten
- einheitliche oder separate Pin-Codes für die Tage der Zeitserien möglich
- flexible definiton der Zeitserien über Cron-Syntax
  - Syntax: <https://github.com/gorhill/cronexpr>
  - Online Editor: <https://www.freeformatter.com/cron-expression-generator-quartz.html>

**Nachteile:**
- Benutzer benötigt evtl. für jeden Zutritt einen separaten Pin-Code
- Änderungen an der Konfiguration müssen auf das Schließsystem übertragen werden


## Beispielkonfigurationen

**jeden Montag 16:00 bis 20:00 von Januar bis März den gleichen Pin-Code verwenden**
```
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

**jeden Montag 16:00 bis 20:00 von Januar bis März einen wechselnden Pin-Code verwenden**
```
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

