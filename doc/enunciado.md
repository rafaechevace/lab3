# Práctica 3

Los objetivos de esta práctica son los siguientes:

- Conocer e implementar una API RESTful sencilla.
- Implementar mecanismos de identificación y autenticación de usuarios.
- Implementar mecanismos de confidencialidad utilizando HTTPS.

## Enunciado

Vamos a crear un prototipo de [base de datos como servicio][1]
que consistirá en crear un servicio HTTP RESTful capaz de almacenar
documentos en formato JSON. La base de datos permitirá diferentes
cuentas de usuarios, así como una gestión básica de los mismos.

Para este prototipo, se asumirá que el servidor funciona en el puerto 5000
en la máquina `myserver.local`, por lo tanto, el _endpoint_ base será
`https://myserver.local:5000`. `myserver.local` deberá resolver a la
IP `127.0.0.1`.

### Almacenamiento de documentos

Cada usuario tendrá su espacio donde podrá alojar sus documentos JSON.
Dicho espacio será accesible a través de `/<usuario>` y
los documentos se podrán crear y recuperar a través de `/<usuario>/<id>`,
donde `<id>` es una cadena que identifica al documento de dicho usuario.
Por ejemplo, `https://myserver.local:5000/luffy/plan-para-ser-el-rey-de-los-piratas`
accederá al documento con identificador `plan-para-ser-el-rey-de-los-piratas`
del usuario `luffy`.

Los documentos se almacenarán como archivos en el propio servidor.
Siguiendo el ejemplo anterior, tendríamos un archivo para dicho documento
en `<ruta_raiz>/luffy/plan-para-ser-el-rey-de-los-piratas.json`,
siendo `<ruta_raiz>` configurable.

### Autenticación de usuarios

Los usuarios del servicio tendrán un nombre de usuario (`username`) y
una contraseña asociada (`password`). Para darse de alta, los usuarios utilizarán
el endpoint `/signup`.

Para poder utilizar el servicio, los usuarios necesitan un _token_ temporal
que el propio servidor emite al utilizar `/login`. Este _token_ expira a
los 5 minutos de su creación, así que los usuarios deberán volver a llamar a
`/login` periódicamente para poder obtener un nuevo _token_. Todas las operaciones
que el usuario ejecute mientras el _token_ no expire serán permitidas.

Un usuario puede tener varios _tokens_ activos al mismo tiempo.

### API

Los objetos que se deben implementar, junto con sus verbos, vienen explicados
a continuación. Si un verbo HTTP no está especificado en un objeto,
dicho método no es aceptado y no se debe permitir su llamada.

### `/version`

- `GET`: devuelve la versión del programa.

  Si todo va bien, devuelve un número de versión del programa en formato `X.Y.Z`.
  Ejemplo:

  ```json
  {"version": "0.1.0"}
  ```

### `/signup`

- `POST`: crear un usuario nuevo. Necesita los siguientes argumentos:

  - `username`: el nombre de usuario.
  - `password`: la contraseña de usuario.

  Si el usuario no puede crearse, la respuesta incluirá un JSON con el motivo del rechazo:

  ```json
  {"reason": "Username already used"}
  ```

  ```json
  {"reason": "Too simple password"}
  ```

### `/login`

- `POST`: autenticar un usuario existente. Necesita los siguientes
  argumentos:

  - `username`: el nombre de usuario.
  - `password`: la contraseña de usuario.

  Si todo va bien, devuelve el _token_ de acceso para el usuario con el
  formato:

  ```json
  {"access_token": "123456789"}
  ```


Todas las operaciones descritas a continuación necesitarán usar un _token_
de autenticación válido para ser atendidas correctamente. El _token_ deberá
enviarse en la **cabecera HTTP** `Authorization` de la petición. El
valor de la cabecera deberá ser `token <token-del-usuario>`. Ejemplo:

```
Authorization: token <token-del-usuario>
```

### `/<string:username>/<string:doc_id>`

- `GET`: obtiene el contenido del documento `doc_id` del usuario `username`.

  Si todo va bien, devuelve el contenido íntegro del documento en formato JSON.

- `POST`: crea un nuevo documento con identificador `doc_id` en el usuario
  `username`. Necesita los siguientes argumentos:

  - `doc_content`: el contenido del documento a crear en formato JSON.

  Si, por ejemplo, quisiéramos subir el siguiente documento:

  ```json
  {
    "1": "Encontrar un espadachín para la tripulación",
    "2": "Encontrar una navegante"
  }
  ```

  La petición debería ser:

  ```json
  {
    "doc_content": {
      "1": "Encontrar un espadachín para la tripulación",
      "2": "Encontrar una navegante"
    }
  }
  ```

  **Recordad que un documento JSON puede ser únicamente un _array_ o incluso un valor suelto.**

  Si todo va bien, devolverá el número de bytes escritos en disco con el formato:

  ```json
  {"size": <total_bytes>}
  ```

  Recuerda que sólo debes almacenar el documento, no el JSON que se envía.
  Esto afectará al número total de bytes que se escriban.

- `PUT`: actualiza el contenido del documento `doc_id` del usuario `username`.
  Necesita los siguientes argumentos:

  - `doc_content`: nuevo contenido del documento en formato JSON
    (ver la descripción del método `POST`).

  Si todo va bien, devuelve el número de bytes escritos en disco con el formato:

  ```json
  {"size": <total_bytes>}
  ```

- `DELETE`: borra el documento `doc_id` del usuario `username`.

  Si todo va bien, devuelve una respuesta vacía:

  ```json
  {}
  ```

### `/<string:username>/_all_docs`

- `GET`: obtiene el contenido de todos los documentos del usuario `username`.

  Si todo va bien, se devolverá un objeto en formato JSON con la siguiente estructura:

  ```json
  {
    "nombre_documento1": {
	  ...contenido de nombre_documento1...
	},
    "nombre_documento2": {
	  ...contenido de nombre_documento2...
	},
	...
  }
  ```

## ¿Qué se pide?

1. Implementar el servicio en el lenguaje Go. Durante las sesiones de prácticas
  se propondrá algún framework recomendado con el que se puede implementar el
  servicio.

   - Para cada argumento de entrada, se debe especificar y comprobar
	   convenientemente su validez y aceptación. Es importante estar
	   seguro que los valores dados por el usuario son **correctos y
	   seguros**.

   - Para cada situación de error, se debe devolver un código de error
     HTTP que sea apropiado para el error en cuestión. Por ejemplo, si
     un documento no existe, el código `404` es apropiado, mientras
     que el código `403` no lo es.

   - Se recomienda utilizar un mecanismo de autenticación de usuario
     visto en clase (como ficheros `shadow`, uso de _salt_ y _digest_
     de la contraseña...).

   - En caso de reinicio del servicio, no es necesario que los _tokens_
     sean persistentes. No pasa nada si el programa del usuario tiene que
     volver a pedir el _token_.

   - Evidentemente, los documentos y la autenticación de los usuarios
     tienen que ser persistentes.

   - La autenticación con el _token_ debe hacerse en cada una de las
     peticiones que lo necesiten. La generación del _token_ debe ser lo
     suficientemente buena como para que no haya colisiones entre
     nuevos _tokens_.

1. `README` y documentación: incluir un `README.md` (o `README.pdf`)
  que explique brevemente el proyecto, cómo compilarlo, instalarlo y ejecutarlo.
  Debe contener las instrucciones necesarias para que cualquier profesional,
  sin necesidad de conocer cómo está implementado, pueda ejecutarlo localmente
  sin problemas.

   - En el documento debe explicarse cómo obtener los certificados SSL para desplegar
    el servicio y cómo hacer que un cliente estándar en GNU/Linux [2]
    (por ejemplo, `curl`) pueda acceder al servicio sin problemas de certificados.

   - En este documento se deben también exponer aspectos que aseguren
     la disponibilidad del servicio como la limitación de cuota de
     espacio por usuario, o limitación de peticiones de usuarios por
     un intervalo de tiempo. Estas medidas no son necesarias
     implementarlas pero sí documentarlas para trabajo futuro.

1. Pueden incluirse scripts o herramientas para generar los certificados e instalarlos.

Las prácticas entregadas se probarán con `curl` y la biblioteca de Python `requests`
para su evaluación mediante pruebas automáticas, por lo que es muy importante
implementar la API tal y como se especifica. Además, la evaluación intentará
aprovechar fallos de seguridad inherentes a la especificación de esta práctica.

## ¿Qué se valora?

1. Implementación de pruebas unitarias y de integración.
1. Chequear y comprobar parámetros de entrada convenientemente y
   anticipar posibles fallos de seguridad que puedan aprovechar los
   usuarios.
1. Crear una imagen de contenedor para el servicio.


[1]: https://en.wikipedia.org/wiki/Cloud_database
[2]: https://wiki.debian.org/Firefox/PrivateCertificateAuthority
