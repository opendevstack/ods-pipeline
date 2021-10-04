# Design Documents

Render (including images) using:

```
asciidoctor-pdf -a allow-uri-read software-architecture.adoc
```

## Images

Best to include like this:

```
image::http://www.plantuml.com/plantuml/proxy?cache=no&src=https://raw.githubusercontent.com/opendevstack/ods-pipeline/master/docs/trigger_architecture.puml[Architecture]
```

That way the images are rendered in GitHub. To render them in PDF, ensure you set `-a allow-uri-read`.
