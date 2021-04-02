# Email

## URL Format
*smtp://__`username`__:__`password`__@__`host`__:__`port`__/?fromAddress=__`fromAddress`__&toAddresses=__`recipient1`__[,__`recipient2`__,...]*

## Additional props

*Can be either supplied using the params argument, or through the URL using ?key=value&key=value etc.*.

* __FromAddress__ - e-mail address that the mail are sent from (**Required**)  
 Aliases: `from`  

* __ToAddresses__ - list of recipient e-mails separated by `,` (**Required**)  
 Aliases: `to`  

* __Auth__ - SMTP authentication method  
 Default: `Unknown`  
 Possible values: `None`, `Plain`, `CRAMMD5`, `Unknown`, `OAuth2`  

* __Encryption__ - Encryption method  
 Default: `Auto`  
 Possible values: `None`, `ExplicitTLS`, `ImplicitTLS`, `Auto`  

* __FromName__ - name of the sender  
 Default: *empty*

* __Password__ - authentication password or hash  
 Default: *empty*  

* __Port__ - SMTP server port, common ones are 25, 465, 587 or 2525  
 Default: `25`  

* __Subject__ - the subject of the sent mail  
 Default: `Shoutrrr Notification`  
 Aliases: `title`  

* __UseHTML__ - whether the message being sent is in HTML  
 Default: `No`  

* __UseStartTLS__ - attempt to use SMTP StartTLS encryption  
 Default: `Yes`  

