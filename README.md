# CLIENTE PARA CONSUMIR API DE ENVIO DE ARQUIVOS

## CRIANDO API

Para criar a api pode clonar o projeto: flieserver e usa-lo como referência ou mesmo como sua propria api.

## CONFIGURANDO

Para realizar as conexões deve alterar as seguintes informações:

-URL exemplo: http://api.domain.com/v1/user/create 
-TOKEN retornado pela api esse é o token gerado no momento em que foi criado o usuário

> O cliente usa esse token para realizar as conexões posteriores.

## ENVIANDO ARQUIVOS


```bash
$ cd clientfilego/
$ go install
$ clientfilego path user password
```




