# d.velop App-Vorlage für Go

Diese Vorlage enthält alles, was Sie benötigen, um eine App für d.velop cloud zu entwickeln.

Um alle Aspekte der App-Entwicklung zu demonstrieren, wird ein hypothetischer, aber nicht trivialer Anwendungsfall
*eines Mitarbeiters, der Urlaub beantragt*, implementiert.

## Erste Schritte

Klonen Sie dieses Repository und folgen Sie den [Build-Anweisungen](#build), um die Beispiel-App zum Laufen zu bringen.
Passen Sie anschließend den Code an Ihr eigenes Geschäftsproblem/Ihre App an.

### Voraussetzungen

Für den Build- und Deployment-Prozess der App wird ein Linux-Docker-Container verwendet.
Daher benötigen Sie auf Ihrem lokalen Entwicklungssystem neben Docker (verwenden Sie eine aktuelle Version) nur einen Git-Client
sowie einen Editor oder eine IDE für Go.

In der Regel erfordert die IDE ein lokal installiertes [Go](https://golang.org/). Verwenden Sie mindestens Version 1.11, da dieses
Projekt [Go Modules](https://github.com/golang/go/wiki/Modules) verwendet.

### Build

Starten Sie den Build mit

```
docker-build build
```

Dabei wird eine eigenständige Web-Anwendung `dist/<appname>app.exe` erstellt, die zum Ausführen und Testen Ihrer App
als lokaler Prozess auf Ihrem Entwicklungssystem verwendet werden kann, sowie ein Deployment-Paket für AWS Lambda `dist/lambda`,
das für den Produktions-Deployment Ihrer App in d.velop cloud verwendet werden sollte.

## App lokal ausführen und testen

Starten Sie einfach `dist/<appname>app.exe`, um Ihre App in einer lokalen Entwicklungsumgebung auszuführen und zu testen.
Bitte beachten Sie, dass einige Funktionen wie die Authentifizierung,
die die Anwesenheit weiterer Apps (z. B. IdentityProviderApp) erfordern,
nicht funktionieren werden, da diese Apps auf Ihrem lokalen System nicht verfügbar sind.

## App umbenennen

Sie sollten den Namen der App ändern, sodass er das Geschäftsproblem widerspiegelt, das Sie lösen möchten.

Jeder App-Name in d.velop cloud muss eindeutig sein. Um dies zu gewährleisten, wählt jeder Anbieter/jedes Unternehmen
ein eindeutiges Anbieter-Präfix, das als Namensraum für die Apps dieses Anbieters dient.
Das Präfix kann während des Registrierungsprozesses in d.velop cloud ausgewählt werden.
Wenn Sie ein Anbieter-Präfix wählen, das Ihrem Firmennamen oder einer Abkürzung davon entspricht,
ist es sehr wahrscheinlich, dass es bei der späteren Registrierung Ihrer App in d.velop cloud verfügbar ist.

Wenn Ihr Unternehmen beispielsweise *Super Duper Software GmbH* heißt und die Domäne Ihrer App
*Mitarbeiter beantragen Urlaub* ist, sollte Ihre App
`superdupergmbh-urlaubsantrag`App heißen. Beachten Sie, dass das Suffix `App` in den Konfigurationsdateien nicht verwendet wird.

Apps, die zur d.velop cloud-Kernplattform gehören, haben kein Anbieter-Präfix.

Verwenden Sie das `rename`-Ziel, um Ihre App umzubenennen:

```
docker-build rename NAME=NEUER_APP_NAME
```

Außerdem möchten Sie möglicherweise die folgenden Werte manuell anpassen:

1.  Ändern Sie das `DOMAIN_SUFFIX` zu einer Domain, die Ihnen gehört, z. B. `ihrunternehmen.de`
2.  `go.mod`: Ändern Sie den Modulnamen von `github.com/d-velop/dvelop-app-template-go` zu etwas wie `github.com/<ihrunternehmen>/<appname>`.
    Dies erfordert leider eine Änderung des Importpfads in vielen Go-Dateien.
    Die Funktion „Ersetzen" Ihrer IDE sollte dabei helfen.


**Bitte schließen Sie mindestens Schritt 1 und Schritt 2 ab, bevor Sie Ihre App [deployen](#deployment), da die Namen vieler
AWS-Ressourcen von `APP_NAME` und `DOMAIN_SUFFIX` abgeleitet werden. Eine nachträgliche Änderung erfordert ein
erneutes Deployment der AWS-Ressourcen, was einige Zeit in Anspruch nimmt.**

## Deployment

**Bitte lesen Sie [App umbenennen](#app-umbenennen), bevor Sie mit dem Deployment fortfahren.**

Sie benötigen ein AWS-Konto, um Ihre App zu deployen. Zum Zeitpunkt der Erstellung dieses Dokuments sind einige AWS-Dienste
für eine begrenzte Zeit und Auslastung kostenlos nutzbar.
Prüfen Sie das [Free Tier](https://aws.amazon.com/free/)-Angebot von AWS für die aktuellen Konditionen.

Erstellen Sie manuell einen IAM-Benutzer mit
den entsprechenden Rechten zum Erstellen der in Ihrer Terraform-Konfiguration definierten AWS-Ressourcen.
Sie könnten mit einem Benutzer beginnen, der die Richtlinie `arn:aws:iam::aws:policy/AdministratorAccess` hat,
aber Sie **sollten die Rechte dieses IAM-Benutzers auf ein Minimum beschränken, sobald Sie in Produktion gehen**.

Konfigurieren Sie die AWS-Anmeldedaten des erstellten IAM-Benutzers mit einer der unter
[Configuring the AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) beschriebenen Methoden.
Setzen Sie beispielsweise die Umgebungsvariablen `AWS_ACCESS_KEY_ID` und `AWS_SECRET_ACCESS_KEY`.

**Windows**

```
SET AWS_ACCESS_KEY_ID=<IHRE-ACCESS-KEY-ID>
SET AWS_SECRET_ACCESS_KEY=<IHR-SECRET-ACCESS-KEY>
```

**Linux**

```
export AWS_ACCESS_KEY_ID=<IHRE-ACCESS-KEY-ID>
export AWS_SECRET_ACCESS_KEY=<IHR-SECRET-ACCESS-KEY>
```

Deployen Sie die Lambda-Funktion und alle anderen AWS-Ressourcen wie AWS API Gateway.

```
docker-build deploy
```

Der Build-Container verwendet [Terraform](https://www.terraform.io/) zur Verwaltung der AWS-Ressourcen und zum Deployment
Ihrer Lambda-Funktion. Dieses Tool implementiert einen Soll-Zustands-Mechanismus, d. h. die erste Ausführung dauert einige Zeit,
um alle erforderlichen AWS-Ressourcen bereitzustellen. Nachfolgende Ausführungen deployen nur die Differenz zwischen dem Soll-Zustand
(z. B. die neue Version Ihrer Lambda-Funktion) und dem bereits deployten Zustand (andere AWS-Ressourcen, die sich zwischen
Deployments nicht ändern) und sind deutlich schneller.

### Endpunkt testen

Die Endpunkt-URLs werden am Ende des Deployments protokolliert. Rufen Sie diese einfach in einem Browser auf, um Ihre App zu testen.

```
Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

endpoint = [
    https://xxxxxxxxxx.execute-api.eu-central-1.amazonaws.com/prod/vacationprocess/,
    https://xxxxxxxxxx.execute-api.eu-central-1.amazonaws.com/dev/vacationprocess/
]

```

Um den aktuellen Deployment-Status anzuzeigen, können Sie jederzeit

```
docker-build show
```

aufrufen, ohne Ihr Deployment zu verändern.

### Deployment einer neuen App-Version

Folgen Sie einfach den [Deployment](#deployment)-Schritten. Ein neues Deployment-Paket für die Lambda-Funktion wird automatisch erstellt.

### Zusätzliche AWS-Ressourcen

Die Terraform-Deployment-Konfiguration enthält 2 zusätzliche Module, die standardmäßig deaktiviert sind.
Kommentieren Sie einfach die entsprechenden Zeilen in `/terraform/main.tf` aus, um sie zu verwenden, aber **stellen Sie sicher, dass die DNS-Auflösung
für Ihre gehostete Zone funktioniert, bevor Sie diese Module verwenden**. Lesen Sie die Kommentare in der Terraform-Datei.

#### asset_cdn
Dieses Modul verwendet *AWS CloudFront* als CDN für Ihre statischen Assets. Außerdem können Sie eine
benutzerdefinierte Domain für Ihre Assets anstelle der S3-URL definieren. Ihr Deployment funktioniert auch ohne dieses Modul einwandfrei.

#### api_custom_domain
Dieses Modul ermöglicht es Ihnen, eine benutzerdefinierte Domain für Ihre App-Endpunkte zu definieren. Ein benutzerdefinierter Domainname ist
erforderlich, sobald Sie Ihre App im d.velop cloud Center registrieren, da der Basispfad Ihrer App mit dem Namen Ihrer App beginnen muss.
Anstelle der Standard-Endpunkte

```
    https://xxxxxxxxxx.execute-api.eu-central-1.amazonaws.com/prod/vacationprocess/
    https://xxxxxxxxxx.execute-api.eu-central-1.amazonaws.com/dev/vacationprocess/
```
deren Basispfade mit `/prod` oder `/dev` beginnen, benötigen Sie Endpunkte wie

```
    https://vacationprocess.xyzdomainde./vacationprocess
    https://dev.vacationprocess.xyzdomainde./vacationprocess
```
die von diesem Modul bereitgestellt werden.

## Projektstruktur

Die vorgestellte Struktur ist keineswegs verbindlich für d.velop cloud-Apps und ist stark meinungsbasiert.
Sie können die Struktur gerne ändern, wenn sie Ihren Anforderungen nicht entspricht.
Andererseits kostet es viel Zeit, eine logische und nützliche Struktur für Apps zu entwickeln, und wir sind ziemlich sicher,
dass diese Struktur zumindest ein guter Ausgangspunkt ist.
Daher empfehlen wir, sie zu verwenden und sich damit vertraut zu machen, damit Sie Ihre Zeit nicht verschwenden
und sofort mit der Implementierung einer Lösung für Ihr Geschäftsproblem beginnen können.

### Go-Verzeichnisse

#### `/cmd`

Enthält die Hauptanwendungen für dieses Projekt. Das ist die eigenständige Webanwendung `/cmd/app`,
die auf Ihrem lokalen Rechner ausgeführt werden kann, und die Lambda-Funktion `/cmd/lambda` für AWS.

Legen Sie nicht viel Code im Anwendungsverzeichnis ab. Legen Sie diesen Code im Verzeichnis `/domain` ab.

Es ist üblich, eine kleine `main`-Funktion zu haben, die im Wesentlichen die Abhängigkeiten verknüpft und ansonsten
vollständig auf den Code aus dem Verzeichnis `/domain` setzt.

#### `/domain`

Enthält den Großteil des Codes für diese App.

Die Struktur folgt den Prinzipien der
[Clean Architecture](http://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) oder
[Hexagonal Architecture](https://alistair.cockburn.us/coming-soon/)
und trennt den Kern der Domäne von externen Frameworks, der Datenbank und der Benutzeroberfläche.

Das Stammverzeichnis enthält das Herzstück der Domäne und hat keine Abhängigkeiten zu externen „Dingen" wie HTTP oder Datenbanken.

##### `/domain/<Anwendungsfall>`

Jeder Anwendungsfall der Domäne hat sein eigenes Unterverzeichnis, das nach dem Anwendungsfall benannt ist. Sie sollten also
die Geschäftsdomäne einer App, die Sie noch nie gesehen haben, verstehen können, indem Sie den Domain-Ordner öffnen und sich
die Verzeichnisnamen ansehen.

Die Anwendungsfälle haben ebenfalls keine Abhängigkeiten zu externen „Dingen".

##### `/domain/mock`

Enthält Test-Mocks, die für mehr als einen Anwendungsfall relevant sind.

##### `/domain/plugins`

Enthält die Abhängigkeiten zu externen „Dingen" wie einer Datenbank oder dem Aufrufkanal, z. B. HTTP.
Die Idee besteht darin, diese externen „Dinge" als Plugins für die Domäne zu behandeln, um
die Domäne einfach, verständlich und separat testbar zu halten. Nicht zuletzt können Sie externe
Abhängigkeiten wie die Datenbank später ändern, ohne die gesamte App neu zu schreiben, da der relevante Code
nicht über die gesamte Codebasis verteilt ist.

Vielleicht möchten Sie auch
[Clean Architecture](http://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
oder [Hexagonal Architecture](https://alistair.cockburn.us/coming-soon/) lesen.

### `/buildcontainer`

Enthält das `Dockerfile` für den Build-Container. Es wird in einem separaten Verzeichnis aufbewahrt, um
den Build-Kontext klein zu halten, damit das Image so schnell wie möglich erstellt werden kann.

### `/terraform`

Enthält die Terraform-Dateien.

### `/web`

Enthält das Web-Frontend.

Das Frontend-Tooling wird auf ein Minimum reduziert, um das gesamte Projekt so einfach wie möglich zu halten.
Außerdem gibt es hunderte mögliche Kombinationen von Frameworks und Build-Tools,
die für das Frontend verwendet werden können. Daher hat jeder Entwickler seine eigenen Präferenzen bezüglich des Toolings.

Verwenden Sie Ihre bevorzugten Tools für das Frontend und ändern Sie den Ordner `web` entsprechend.
Vergessen Sie nicht, die Aufgabe `deploy-assets` im `makefile` und die go:generate-Befehle in `/domain/plugins/gui/` zu ändern.

Es ist wahrscheinlich, dass wir in Zukunft Web-Projekte mit verschiedenen Tools und Frameworks bereitstellen werden,
die den Ordner `web` ersetzen können.

## Go Modules

Dieses Projekt verwendet [Go Modules](https://golang.org/doc/go1.11#modules). Das bedeutet, Sie benötigen mindestens Go 1.11, wenn Sie
dieses Projekt außerhalb des Build-Containers kompilieren möchten. Das bedeutet auch, dass Ihr Projekt **nicht in GOPATH/src liegen darf**
(vgl. [Preliminary module support](https://golang.org/cmd/go/#hdr-Preliminary_module_support)) und die **Abhängigkeiten
nicht in die Quellcodeverwaltung eingecheckt werden dürfen**.

### IDE-Unterstützung für Go Modules

In einigen IDEs, wie JetBrains GoLand, muss die Go-Modules-Unterstützung explizit aktiviert werden, um IntelliSense zu erhalten.
* Einstellungen > Go > Go Module (vgo) - Go Modules (vgo)-Integration aktivieren

## Build-Mechanismus
Ein Linux-Docker-Container wird zum Erstellen und Deployen der Software verwendet. Dies hat den Vorteil, dass der Build
nicht von bestimmten Tools oder Tool-Versionen abhängt, die auf dem lokalen Entwicklungsrechner oder
Build-Server installiert sein müssen.

Während des Builds wird das gesamte Anwendungsverzeichnis in den Docker-Container eingebunden. Die Build-Ziele sind
im `Makefile` implementiert.

Zwei Wrapper (`docker-build.bat` und `docker-build.sh`) werden bereitgestellt, damit Sie sich den
relativ langen Docker-Befehl nicht merken müssen.
Außerdem bieten diese Wrapper eine kleine Hilfsfunktion, um alle in der `environment`-Datei aufgelisteten Umgebungsvariablen
vom Docker-Host (d. h. Ihrem Entwicklungsrechner oder Build-Server)
an den Build-Container weiterzureichen.

## Mitwirken

Bitte lesen Sie [CONTRIBUTING.md](CONTRIBUTING.md) für Details zu unserem Verhaltenskodex und den Prozess zum Einreichen von Pull Requests.

## Lizenz

Bitte lesen Sie [LICENSE](LICENSE) für Lizenzinformationen.

## Danksagungen

Dank an die folgenden Projekte für die Inspiration:

* [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
* [How Do You Structure Your Go Apps](https://github.com/katzien/go-structure-examples)
* [GoDDD](https://github.com/marcusolsson/goddd)
* [Starting an Open Source Project](https://opensource.guide/starting-a-project/)
* [README template](https://gist.github.com/PurpleBooth/109311bb0361f32d87a2)
* [CONTRIBUTING template](https://github.com/nayafia/contributing-template/blob/master/CONTRIBUTING-template.md)
