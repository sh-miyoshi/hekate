# API specification

## Overview

API specification is writen in OpenAPI format.  
So, you can view `app.yaml` in [swagger editor](http://editor.swagger.io/).

## How to generate HTML

```bash
npm install -g redoc-cli
redoc-cli bundle api.yaml -o api.html
```
