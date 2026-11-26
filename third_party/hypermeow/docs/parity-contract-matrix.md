# Matriz de contratos: 303 operações ZuckZapGo e equivalentes WPPConnect

Data: 10/11/2026. Esta é uma comparação estática das fontes fixadas, não uma certificação de paridade funcional. O [JSON estruturado](parity-contract-matrix.compact.json) contém uma entrada para cada uma das 303 operações, schemas resumidos de entrada/resposta, candidatos Go/Node, diferenças de campos, fonte/linha e critérios de aceite.

## Baselines e método

- ZuckZapGo: imagem `setupautomatizado/zuckzapgo-private:v2.4.2`, manifest AMD64 `sha256:f073c6e57ed95b6dd9a75574c2e60bfef0b302243df1d07fc11260caa6381c1a`; spec extraído da camada final, SHA-256 `2096ea73209320fb4489f3751fd563fd47bf9c00a6acb6191314a1cd814db7a7`. São 270 paths e 303 operações.
- Go: [`052208a8818c942241d867d385d59c88df3792f9`](https://github.com/wppconnect-team/wppconnect-server-go/blob/052208a8818c942241d867d385d59c88df3792f9/docs/openapi.json), 135 paths/140 operações.
- Node: [`37b7bc1a662b5262b8276597701f38c4dbcad58a`](https://github.com/wppconnect-team/wppconnect-server/blob/37b7bc1a662b5262b8276597701f38c4dbcad58a/src/swagger.json), 157 paths/158 operações.
- As 81 coincidências método/path entre Go e Node constam separadamente no JSON. Não representam 81 contratos equivalentes.
- Mapeamento semântico explícito, sem fuzzy matching. Cada candidato foi validado como método/path presente no spec de destino. Handlers base foram consultados para limitações de grupos, mensagens, mídia, histórico, administração e canais.
- Parser YAML 2.8.2 obtido do registry npm, tarball conferido por SHA-512. A extração lê fontes fixadas via Git e verifica hash e contagem. O JSON público deduplica valores em `$defs` com referências locais `$ref`; sua expansão foi comparada por igualdade profunda com o inventário completo, preservando todas as 303 operações, critérios e limitações. A chave `compactEncoding` documenta o formato e o hash do inventário original.

O JSON usa `functionalParityValidated:false` em **todos** os destinos: nenhum teste ao vivo foi realizado. `core` significa que existe a operação central em código/contrato e pode receber um adapter; não significa que o request Zuck pode ser enviado diretamente. `partial` preserva somente parte do escopo; `composition` exige orquestração adicional; `related` é apenas infraestrutura semelhante. `unsupported_501` e `experimental` não entram como cobertura estável. `missing` é ausência de equivalente exposto identificado no baseline, não ausência absoluta de capacidade do SDK.

| Classificação | Go | Node |
| --- | ---: | ---: |
| núcleo adaptável | 56 | 66 |
| parcial | 50 | 37 |
| composição | 5 | 4 |
| somente relacionado | 4 | 5 |
| experimental | 2 | 0 |
| 501 explícito | 1 | 0 |
| sem equivalente exposto | 184 | 190 |
| fora do escopo | 1 | 1 |

Não somar essas colunas para anunciar percentual de paridade: uma operação com muitos efeitos pode ter várias rotas candidatas, e uma rota pode atender apenas uma parte de muitas operações.

## Diferenças de payload verificáveis

O JSON resolve refs locais e registra campos, tipos, required, enums, limites, parâmetros de path/query, content types, códigos de resposta e esquemas de autenticação documentados. `sourceOnly`, `targetOnly` e `sharedNames` comparam somente nomes de campos JSON no primeiro nível. Campo com nome igual pode ter significado diferente; campo com nome diferente pode ser adaptável. Enums e limites ausentes no schema não são inferidos a partir de silêncio da documentação. Descrições e exemplos extensos não foram copiados.

Exemplos relevantes:

| Jornada | Diferença confirmada | Implicação |
| --- | --- | --- |
| Texto | Zuck Phone/Body e opções Id/Presence/Duration/LinkPreview/ContextInfo; Go phone/message e quotedMessageId retido | Adapter explícito; não prometer todas as opções via rename |
| Imagem | Zuck Image pode ser URL; Go requer base64/data URL/multipart e limita 16 MiB | Não fazer fetch arbitrário implícito nem ignorar opção |
| Grupos | locked usa groupjid/locked; Go groupId/property=restrict/value | Mapear false corretamente e confirmar estado remoto |
| Reações/reply | Go exige referência observada na mesma sessão/chat | IDs externos ou histórico perdido não podem ser aceitos como equivalentes |
| Labels | Definições, associação de chat e de mensagem são operações diferentes | Não mapear todo label/manage para uma única chamada genérica |
| Histórico | Go usa arquivo observado parcial e cursor próprio | Não retornar histórico parcial sob promessa de full-history |
| Webhook | Zuck CRUD por instância; Go base só WEBHOOK_URL do processo | Configurar pool/tenant e persistência antes de anunciar suporte |
| Calls | Go/Node base oferecem rejeição; demais controles/mídia não têm equivalente identificado | Não declarar VoIP privado implementado |

Fontes de comportamento: [Go router](https://github.com/wppconnect-team/wppconnect-server-go/blob/052208a8818c942241d867d385d59c88df3792f9/internal/httpapi/router.go), [Go groups](https://github.com/wppconnect-team/wppconnect-server-go/blob/052208a8818c942241d867d385d59c88df3792f9/internal/httpapi/groups.go), [Go media](https://github.com/wppconnect-team/wppconnect-server-go/blob/052208a8818c942241d867d385d59c88df3792f9/internal/httpapi/media.go), [Go channels](https://github.com/wppconnect-team/wppconnect-server-go/blob/052208a8818c942241d867d385d59c88df3792f9/internal/httpapi/channels.go), [Node routes](https://github.com/wppconnect-team/wppconnect-server/blob/37b7bc1a662b5262b8276597701f38c4dbcad58a/src/routes/index.ts). O contrato da imagem é evidência local identificada por digest; a [página DockerHub](https://hub.docker.com/r/setupautomatizado/zuckzapgo-private) e a [release v2.4.2](https://github.com/guilhermejansen/use-zuckzapgo/releases/tag/v2.4.2) identificam o artefato, sem substituir a inspeção exata.

## Jornadas e aceite

| ID | Jornada | Prioridade | Aceite mínimo |
| --- | --- | --- | --- |
| J01 | Sessão e pairing | P0 | Criar duas sessões isoladas; QR/phone code somente para a sessão autorizada. Conectar, reiniciar e restaurar sem perder credenciais; distinguir disconnect de logout. Responder corretamente a QR expirado, sessão fechada e pairing recusado. |
| J02 | Webhooks e configuração por instância | P0 | Criar pela Cloud usando pool com webhookUrl/webhook e receber evento no destino correto. Persistir configuração no restart; não entregar evento de outro tenant. Testar assinatura, deduplicação, timeout/429/503, URL inválida e remoção de destino. |
| J03 | Mensagens de texto e ações | P1 | Texto/reply/edit/delete/reaction confirmados em dois dispositivos autorizados, com IDs persistidos. Verificar PN/LID/grupo, autoria, janela de edição e erro por destinatário. Timeout/retry não pode reenviar silenciosamente; distinguir ACK de entrega. |
| J04 | Upload e download de mídia | P1 | Imagem, áudio/PTT, documento, MP4 e WebP enviados/baixados com MIME e hash corretos. Rejeitar excesso de tamanho, arquivo inválido e referência de outro tenant. Testar mídia expirada, sessão offline, conversão quando suportada e limpeza de temporários. |
| J05 | Interativos e enquetes | P3 | Validar UI e respostas em Android/iOS/Web e contas relevantes; ACK sozinho não aprova. Conferir limites, escolha de opções, voto/retirada e referências criptográficas corretas. Desabilitado/sem suporte devolve capacidade explícita e não sucesso simulado. |
| J06 | Chats, histórico e sincronização | P1 | Paginação sem perda/duplicação, ordenação e sinalização de arquivo parcial. Ações de chat só atualizam índice depois do ACK; validar restart. Não afirmar histórico completo quando faltam mensagens ou permissões de sync. |
| J07 | Contatos, identidade e perfil | P2 | Normalizar PN/LID/username sem criar identidades duplicadas. Conferir permissões, dados ausentes, perfil próprio versus terceiro e persistência. Campos privados e tokens não aparecem em resposta/logs; false/vazio/ausente têm semântica testada. |
| J08 | Grupos e participantes | P2 | Criar grupo, convite, entrar/sair e ações com admin e membro comum. Preservar resultados por participante e falhas parciais; testar propriedades false/zero. Conferir estado remoto após mutação, não apenas HTTP 200. |
| J09 | Comunidades | P3 | Distinguir parent, subgrupo e grupo de anúncios; verificar vínculo real. Testar permissão, pedido pendente, aprovação/rejeição e anúncio no destino correto. Não contar envio genérico como roteamento de comunidade pronto. |
| J10 | Newsletters | P2 | Criar/listar/follow/mute/send e consultar como owner/subscriber com permissões reais. Conferir mídia, metadados, IDs de servidor, paginação e expiração de live updates. Sem suporte administrativo deve ser explícito, sem substituir por operação diferente. |
| J11 | Catálogo e perfil business | P3 | Ler/criar/editar produtos/coleções com conta business autorizada e conferir estado remoto. Validar moeda, precisão, imagem, visibilidade false, paginação e permissão por catálogo. Testar idempotência e resultados parciais sem dados de outro negócio. |
| J12 | Pedidos e pagamentos | P3 | Relacionar orderID/messageID/referenceID e preservar valores monetários exatamente. Testar transições permitidas e reprocessamento idempotente. Mensagem de cobrança/status não deve ser declarada confirmação externa de pagamento. |
| J13 | Chamadas e mídia em tempo real | P4 | Rejeição é validada separadamente; dial/answer/record só aprovam com chamada real autorizada. Para mídia: áudio bidirecional, codec, perda/reconexão, vídeo e teardown sem vazamento. Testar autenticação/scopes, NAT/UDP e limites; endpoint documentado não prova engine implementada. |
| J14 | Operação, eventos, arquivo e DLQ | P1 | Derrubar destino, observar backlog/failed, reiniciar e recuperar sem perder registro persistido. Replay por ID preserva deduplicação; lote/prune requer filtros, ownership e auditoria. Health/capabilities/metrics informam estado real; nenhuma credencial nos eventos. |
| J15 | Broker/S3 por sessão | P2 | Testar dois tenants, broker/S3 reais, autorização e indisponibilidade seguida de recuperação. Conferir confirmações, retenção, quota e configuração persistida/redigida. Configuração global do processo não conta como CRUD por sessão. |
| J16 | Connector, passkey e import | P4 | Token curto só permite pairing, expira e respeita origem; token permanente não vai ao browser. Testar desafio passkey e import com material autorizado, validação e redaction. Sem implementação explícita, retornar indisponível; não reutilizar QR administrativo como connector seguro. |
| J17 | Privacidade e status | P2 | Testar permissões/audiência, block/unblock, timer e profile/status sem vazamento. Validar leitura depois de escrita, zero/false/vazio e conta sem permissão. Status áudio/imagem/vídeo precisa ser recebido por audiência autorizada. |
| J18 | Rota mobile e proxy | P4 | Proxy comum e rota mobile têm contratos separados; testar fail-closed e versionamento. Nenhum fallback direto quando rota obrigatória falha; segredos não persistidos indevidamente. Conferir reconnect/restart e restrições TCP/UDP para chamadas. |
| J19 | Labels | P2 | Criar/editar/aplicar/remover e conferir sync em outro dispositivo. Distinguir associação chat/mensagem e definição da etiqueta. Preservar tombstones/limites e isolamento no restart. |
| J20 | Licença comercial do produto de referência | fora-do-escopo | Não copiar o mecanismo comercial do concorrente como requisito de paridade WhatsApp. Documentar capacidades e versão próprias do WPPConnect com atribuição correta. |

## Pendências concretas priorizadas

1. **P0: Cloud Go-only e webhook no pool.** Validar seleção em todas as camadas e transportar webhookUrl/webhook até a ativação do pool. O baseline do worker descarta opções ao ativar container, e Go ignora body de start-session. Aceite exige evento real após criação pela API pública e restart, não só teste unitário do helper.
2. **P1: contrato de mensagens, mídia e ações.** Implementar adapter onde há núcleo, erros explícitos para campos sem suporte, idempotência e testes de resultado ambíguo. Cobrir JSON/multipart, limites e isolamento.
3. **P1: operação de eventos/histórico.** Dar visibilidade a backlog e falha, DLQ/replay auditável, paginação e retenção; não representar stream por sessão como stream administrativo global.
4. **P2: capacidades já disponíveis no SDK.** Newsletters views/reactions/update, contatos/business reads, timers e favoritos podem ser entregas pequenas, mas só contam na matriz quando expostas e testadas no servidor.
5. **P2: transportes/S3 configuráveis por tenant.** A infraestrutura global existente não supre CRUD por sessão, filtros, teste de conexão e segredos redigidos.
6. **P3: business/interativos.** Reutilizar HyperMeow público com migração/imports/schema validada; acrescentar handlers e fixtures para catálogo, pedido, flow/carousel. Validar entrega real sem promover legacy experimental a paridade.
7. **P4: calls/connector/passkey/import/mobile route.** Projetos separados com aceite próprio. A imagem contém documentação, não fonte compilável da engine privada; nenhum trabalho aqui implementou áudio/vídeo.

Este arquivo mantém os baselines originais para comparação reproduzível. Alterações locais posteriores de HMAC, Cloud/Manager/worker ou HyperMeow não são retroativamente contadas. Para consultar os contratos de origem, use os links fixados por commit acima. O SDK público de referência é [polymorfa/hypermeow no commit f432198](https://github.com/polymorfa/hypermeow/tree/f43219816d8600ec963371881a1467d22fa7b3e6); sua existência não comprova os endpoints da aplicação nem o conteúdo do fork privado usado na imagem.

## Inventário completo

Nos candidatos abaixo, paths sem prefixo são relativos a `/api/{session}/`. Método é preservado; caminhos administrativos e globais aparecem completos. Cada ID abre a entrada correspondente no JSON por busca textual. A nota e as diferenças de campos por operação estão no JSON para manter esta tabela legível.

| ID | Operação ZuckZapGo | Go | Node | Aceite / prioridade |
| --- | --- | --- | --- | --- |
| Z001 | `POST /session/echo/api` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z002 | `GET /session/echo/api` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z003 | `GET /admin/users` | parcial: `GET /api/admin/sessions` | parcial: `POST /api/show-all-sessions` | J14 / P1 |
| Z004 | `POST /admin/users` | parcial: `POST /api/admin/sessions` | sem equivalente exposto | J14 / P1 |
| Z005 | `GET /admin/instances/status` | parcial: `GET /api/admin/sessions` | parcial: `POST /api/show-all-sessions` | J14 / P1 |
| Z006 | `GET /admin/users/{id}` | parcial: `GET /api/admin/session-settings/{name}` | sem equivalente exposto | J14 / P1 |
| Z007 | `PUT /admin/users/{id}` | parcial: `PUT /api/admin/session-settings/{name}` | sem equivalente exposto | J14 / P1 |
| Z008 | `DELETE /admin/users/{id}` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z009 | `DELETE /admin/users/{id}/full` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z010 | `GET /admin/license/info` | fora do escopo | fora do escopo | J20 / fora-do-escopo |
| Z011 | `GET /newsletter/list` | parcial: `GET newsletters` | sem equivalente exposto | J10 / P2 |
| Z012 | `POST /newsletter/create` | núcleo adaptável: `POST create-newsletter` | núcleo adaptável: `POST newsletter` | J10 / P2 |
| Z013 | `GET /newsletter/subscribed` | parcial: `GET newsletters` | sem equivalente exposto | J10 / P2 |
| Z014 | `POST /newsletter/info` | núcleo adaptável: `POST newsletter-info` | sem equivalente exposto | J10 / P2 |
| Z015 | `POST /newsletter/send` | parcial: `POST send-newsletter-text` | sem equivalente exposto | J10 / P2 |
| Z016 | `POST /newsletter/follow` | núcleo adaptável: `POST follow-newsletter` | sem equivalente exposto | J10 / P2 |
| Z017 | `POST /newsletter/follow-invite` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z018 | `POST /newsletter/unfollow` | núcleo adaptável: `POST unfollow-newsletter` | sem equivalente exposto | J10 / P2 |
| Z019 | `POST /newsletter/unfollow-invite` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z020 | `POST /newsletter/mute` | núcleo adaptável: `POST mute-newsletter` | núcleo adaptável: `POST mute-newsletter/{id}` | J10 / P2 |
| Z021 | `POST /newsletter/messages` | parcial: `POST newsletter-messages` | sem equivalente exposto | J10 / P2 |
| Z022 | `POST /newsletter/reaction` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z023 | `POST /newsletter/mark-viewed` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z024 | `POST /newsletter/subscribe-live-updates` | núcleo adaptável: `POST subscribe-newsletter` | sem equivalente exposto | J10 / P2 |
| Z025 | `POST /newsletter/message-updates` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z026 | `POST /newsletter/accept-tos` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z027 | `POST /newsletter/admin-invite` | sem equivalente exposto | sem equivalente exposto | J10 / P2 |
| Z028 | `GET /health` | parcial: `GET /healthz`; `GET /readyz`; `GET /metrics` | parcial: `GET /healthz`; `GET /metrics` | J14 / P1 |
| Z029 | `GET /webhook` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z030 | `POST /webhook` | sem equivalente exposto | parcial: `POST start-session` | J02 / P0 |
| Z031 | `PUT /webhook` | sem equivalente exposto | parcial: `POST start-session` | J02 / P0 |
| Z032 | `DELETE /webhook` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z033 | `POST /session/connect` | parcial: `POST start-session` | parcial: `POST start-session` | J01 / P0 |
| Z034 | `POST /session/disconnect` | núcleo adaptável: `POST close-session` | núcleo adaptável: `POST close-session` | J01 / P0 |
| Z035 | `POST /session/logout` | núcleo adaptável: `POST logout-session` | núcleo adaptável: `POST logout-session` | J01 / P0 |
| Z036 | `POST /session/refresh` | composição: `POST close-session`; `POST start-session` | composição: `POST close-session`; `POST start-session` | J01 / P0 |
| Z037 | `GET /session/status` | núcleo adaptável: `GET status-session` | núcleo adaptável: `GET status-session` | J01 / P0 |
| Z038 | `POST /session/pairphone` | núcleo adaptável: `POST pair-phone` | sem equivalente exposto | J01 / P0 |
| Z039 | `GET /session/passkey/status` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z040 | `POST /session/passkey/response` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z041 | `POST /session/passkey/confirm` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z042 | `POST /session/passkey/challenge` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z043 | `POST /session/connector/token` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z044 | `POST /connector/connect` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z045 | `POST /connector/disconnect` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z046 | `GET /connector/status` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z047 | `GET /connector/qr` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z048 | `POST /connector/pairphone` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z049 | `GET /connector/passkey/status` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z050 | `POST /connector/passkey/challenge` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z051 | `POST /connector/passkey/response` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z052 | `POST /connector/passkey/confirm` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z053 | `POST /session/import` | sem equivalente exposto | sem equivalente exposto | J16 / P4 |
| Z054 | `GET /session/qr` | núcleo adaptável: `GET qrcode-session` | núcleo adaptável: `GET qrcode-session` | J01 / P0 |
| Z055 | `POST /session/network-route` | sem equivalente exposto | sem equivalente exposto | J18 / P4 |
| Z056 | `GET /session/network-route` | sem equivalente exposto | sem equivalente exposto | J18 / P4 |
| Z057 | `DELETE /session/network-route` | sem equivalente exposto | sem equivalente exposto | J18 / P4 |
| Z058 | `POST /session/proxy` | parcial: `PUT /api/admin/session-settings/{name}` | sem equivalente exposto | J18 / P4 |
| Z059 | `POST /session/s3/config` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z060 | `GET /session/s3/config` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z061 | `DELETE /session/s3/config` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z062 | `POST /session/s3/test` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z063 | `POST /session/rabbitmq/config` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z064 | `GET /session/rabbitmq/config` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z065 | `DELETE /session/rabbitmq/config` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z066 | `POST /session/rabbitmq/test` | sem equivalente exposto | sem equivalente exposto | J15 / P2 |
| Z067 | `POST /session/skipmedia/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z068 | `GET /session/skipmedia/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z069 | `POST /session/skipgroups/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z070 | `GET /session/skipgroups/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z071 | `POST /session/skipnewsletters/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z072 | `GET /session/skipnewsletters/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z073 | `POST /session/skipbroadcasts/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z074 | `GET /session/skipbroadcasts/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z075 | `POST /session/skipownmessages/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z076 | `GET /session/skipownmessages/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z077 | `POST /session/skipcalls/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z078 | `GET /session/skipcalls/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z079 | `POST /session/globaltransports/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z080 | `GET /session/globaltransports/config` | sem equivalente exposto | sem equivalente exposto | J02 / P0 |
| Z081 | `POST /session/presence/config` | somente relacionado: `PUT /api/admin/session-settings/{name}` | sem equivalente exposto | J01 / P0 |
| Z082 | `GET /session/presence/config` | somente relacionado: `GET /api/admin/session-settings/{name}` | sem equivalente exposto | J01 / P0 |
| Z083 | `POST /user/info` | parcial: `GET contact/{phone}`; `GET profile-status/{phone}`; `GET profile-pic/{phone}` | parcial: `GET profile/{phone}`; `GET contact/{phone}` | J07 / P2 |
| Z084 | `GET /user/privacy/settings` | núcleo adaptável: `GET privacy-settings` | sem equivalente exposto | J07 / P2 |
| Z085 | `POST /user/pushname` | núcleo adaptável: `POST set-profile-name` | núcleo adaptável: `POST change-username` | J07 / P2 |
| Z086 | `POST /user/contact/info` | parcial: `GET contact/{phone}` | parcial: `GET contact/{phone}` | J07 / P2 |
| Z087 | `POST /user/contact/add` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z088 | `POST /user/contact/remove` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z089 | `GET /user/status` | composição: `GET profile-status/{phone}` | composição: `GET profile-status/{phone}` | J07 / P2 |
| Z090 | `POST /user/status` | parcial: `POST set-profile-status` | núcleo adaptável: `POST profile-status` | J07 / P2 |
| Z091 | `POST /user/update` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z092 | `POST /user/business/profile` | sem equivalente exposto | somente relacionado: `GET profile/{phone}` | J07 / P2 |
| Z093 | `POST /user/check` | núcleo adaptável: `GET check-number-status/{phone}` | núcleo adaptável: `GET check-number-status/{phone}` | J07 / P2 |
| Z094 | `POST /user/presence` | núcleo adaptável: `POST set-online-presence` | núcleo adaptável: `POST set-online-presence` | J07 / P2 |
| Z095 | `POST /user/avatar` | núcleo adaptável: `GET profile-pic/{phone}` | núcleo adaptável: `GET profile-pic/{phone}` | J07 / P2 |
| Z096 | `POST /user/photo` | núcleo adaptável: `POST set-profile-picture` | núcleo adaptável: `POST set-profile-pic` | J07 / P2 |
| Z097 | `POST /user/photo/remove` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z098 | `GET /user/contacts` | núcleo adaptável: `GET all-contacts` | núcleo adaptável: `GET all-contacts` | J07 / P2 |
| Z099 | `GET /user/lid/get` | parcial: `GET contact/pn-lid/{pnLid}` | parcial: `GET contact/pn-lid/{pnLid}` | J07 / P2 |
| Z100 | `POST /user/lid/get` | parcial: `GET contact/pn-lid/{pnLid}` | parcial: `GET contact/pn-lid/{pnLid}` | J07 / P2 |
| Z101 | `GET /user/lid/from-lid` | parcial: `GET contact/pn-lid/{pnLid}` | parcial: `GET contact/pn-lid/{pnLid}` | J07 / P2 |
| Z102 | `POST /user/lid/from-lid` | parcial: `GET contact/pn-lid/{pnLid}` | parcial: `GET contact/pn-lid/{pnLid}` | J07 / P2 |
| Z103 | `GET /user/lid/mappings` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z104 | `POST /chat/delete` | parcial: `POST delete-message` | núcleo adaptável: `POST delete-message` | J03 / P1 |
| Z105 | `POST /chat/markread` | núcleo adaptável: `POST mark-read` | parcial: `POST send-seen` | J03 / P1 |
| Z106 | `POST /chat/react` | parcial: `POST send-reaction` | núcleo adaptável: `POST react-message` | J03 / P1 |
| Z107 | `POST /chat/send/text` | parcial: `POST send-message` | parcial: `POST send-message` | J03 / P1 |
| Z108 | `POST /chat/send/link` | somente relacionado: `POST send-message` | núcleo adaptável: `POST send-link-preview` | J03 / P1 |
| Z109 | `POST /chat/send/edit` | parcial: `POST edit-message` | núcleo adaptável: `POST edit-message` | J03 / P1 |
| Z110 | `POST /chat/send/image` | núcleo adaptável: `POST send-image` | núcleo adaptável: `POST send-image` | J04 / P1 |
| Z111 | `POST /chat/send/audio` | parcial: `POST send-audio-base64`; `POST send-voice-base64` | parcial: `POST send-file`; `POST send-voice-base64` | J04 / P1 |
| Z112 | `POST /chat/send/document` | núcleo adaptável: `POST send-file` | núcleo adaptável: `POST send-file` | J04 / P1 |
| Z113 | `POST /chat/send/video` | núcleo adaptável: `POST send-video` | parcial: `POST send-file` | J04 / P1 |
| Z114 | `POST /chat/send/ptv` | sem equivalente exposto | sem equivalente exposto | J04 / P1 |
| Z115 | `POST /chat/send/sticker` | núcleo adaptável: `POST send-sticker` | núcleo adaptável: `POST send-sticker` | J04 / P1 |
| Z116 | `POST /chat/send/location` | núcleo adaptável: `POST send-location` | núcleo adaptável: `POST send-location` | J03 / P1 |
| Z117 | `POST /chat/send/contact` | núcleo adaptável: `POST send-contact` | parcial: `POST contact-vcard` | J03 / P1 |
| Z118 | `POST /chat/send/list` | experimental: `POST send-list` | núcleo adaptável: `POST send-list-message` | J05 / P3 |
| Z119 | `POST /chat/send/poll` | núcleo adaptável: `POST send-poll-message` | núcleo adaptável: `POST send-poll-message` | J05 / P3 |
| Z120 | `POST /chat/send/pollvote` | parcial: `POST vote-poll` | sem equivalente exposto | J05 / P3 |
| Z121 | `POST /chat/send/eventresponse` | sem equivalente exposto | sem equivalente exposto | J05 / P3 |
| Z122 | `POST /chat/send/contact-request` | sem equivalente exposto | sem equivalente exposto | J05 / P3 |
| Z123 | `POST /chat/send/event` | sem equivalente exposto | sem equivalente exposto | J05 / P3 |
| Z124 | `POST /chat/send/buttons` | experimental: `POST send-buttons` | núcleo adaptável: `POST send-buttons` | J05 / P3 |
| Z125 | `POST /chat/send/flow` | sem equivalente exposto | sem equivalente exposto | J05 / P3 |
| Z126 | `POST /chat/send/carousel` | 501 explícito: `POST send-carousel` | sem equivalente exposto | J05 / P3 |
| Z127 | `POST /chat/downloadimage` | parcial: `POST download-media` | parcial: `POST download-media` | J04 / P1 |
| Z128 | `POST /chat/downloadvideo` | parcial: `POST download-media` | parcial: `POST download-media` | J04 / P1 |
| Z129 | `POST /chat/downloaddocument` | parcial: `POST download-media` | parcial: `POST download-media` | J04 / P1 |
| Z130 | `POST /chat/downloadaudio` | parcial: `POST download-media` | parcial: `POST download-media` | J04 / P1 |
| Z131 | `POST /chat/downloadsticker` | parcial: `POST download-media` | parcial: `POST download-media` | J04 / P1 |
| Z132 | `POST /chat/presence` | núcleo adaptável: `POST chat-state` | núcleo adaptável: `POST chat-state` | J06 / P1 |
| Z133 | `POST /chat/pin` | núcleo adaptável: `POST pin-chat` | núcleo adaptável: `POST pin-chat` | J06 / P1 |
| Z134 | `POST /chat/archive` | núcleo adaptável: `POST archive-chat` | núcleo adaptável: `POST archive-chat` | J06 / P1 |
| Z135 | `POST /chat/mute` | núcleo adaptável: `POST mute-chat` | núcleo adaptável: `POST send-mute` | J06 / P1 |
| Z136 | `POST /chat/star` | sem equivalente exposto | núcleo adaptável: `POST star-message` | J06 / P1 |
| Z137 | `POST /chat/delete-chat` | parcial: `POST delete-chat` | núcleo adaptável: `POST delete-chat` | J06 / P1 |
| Z138 | `POST /chat/favorite` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z139 | `GET /chat/favorites` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z140 | `PUT /chat/favorites` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z141 | `POST /chat/pin-message` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z142 | `GET /chat/list` | parcial: `POST list-chats`; `GET all-chats` | parcial: `POST list-chats`; `GET all-chats` | J06 / P1 |
| Z143 | `POST /sync/app-state` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z144 | `POST /sync/history-request` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z145 | `POST /sync/full-history` | sem equivalente exposto | sem equivalente exposto | J06 / P1 |
| Z146 | `POST /label/manage` | parcial: `POST edit-label` | parcial: `POST add-new-label`; `PUT delete-label/{id}` | J19 / P2 |
| Z147 | `POST /label/chat` | núcleo adaptável: `POST label-chat` | núcleo adaptável: `POST add-or-remove-label` | J19 / P2 |
| Z148 | `POST /label/message` | núcleo adaptável: `POST label-message` | sem equivalente exposto | J19 / P2 |
| Z149 | `GET /label/list` | núcleo adaptável: `GET labels` | núcleo adaptável: `GET get-all-labels` | J19 / P2 |
| Z150 | `POST /group/create` | núcleo adaptável: `POST create-group` | núcleo adaptável: `POST create-group` | J08 / P2 |
| Z151 | `POST /group/locked` | núcleo adaptável: `POST group-property` | núcleo adaptável: `POST group-property` | J08 / P2 |
| Z152 | `POST /group/ephemeral` | núcleo adaptável: `POST group-property` | núcleo adaptável: `POST group-property` | J08 / P2 |
| Z153 | `POST /group/photo/remove` | sem equivalente exposto | sem equivalente exposto | J08 / P2 |
| Z154 | `GET /group/list` | núcleo adaptável: `GET all-groups` | núcleo adaptável: `GET all-groups` | J08 / P2 |
| Z155 | `GET /group/invitelink` | parcial: `GET group-invite-link/{groupId}`; `GET group-revoke-link/{groupId}` | parcial: `GET group-invite-link/{groupId}`; `GET group-revoke-link/{groupId}` | J08 / P2 |
| Z156 | `GET /group/info` | núcleo adaptável: `GET group-info/{groupId}` | núcleo adaptável: `GET group-info/{groupId}` | J08 / P2 |
| Z157 | `POST /group/photo` | núcleo adaptável: `POST group-pic` | núcleo adaptável: `POST group-pic` | J08 / P2 |
| Z158 | `POST /group/leave` | núcleo adaptável: `POST leave-group` | núcleo adaptável: `POST leave-group` | J08 / P2 |
| Z159 | `POST /group/name` | núcleo adaptável: `POST group-subject` | núcleo adaptável: `POST group-subject` | J08 / P2 |
| Z160 | `POST /group/topic` | núcleo adaptável: `POST group-description` | núcleo adaptável: `POST group-description` | J08 / P2 |
| Z161 | `POST /group/announce` | núcleo adaptável: `POST messages-admins-only` | núcleo adaptável: `POST messages-admins-only` | J08 / P2 |
| Z162 | `POST /group/memberaddmode` | sem equivalente exposto | somente relacionado: `POST group-property` | J08 / P2 |
| Z163 | `POST /group/joinapproval` | sem equivalente exposto | somente relacionado: `POST group-property` | J08 / P2 |
| Z164 | `POST /group/join` | núcleo adaptável: `POST join-code` | núcleo adaptável: `POST join-code` | J08 / P2 |
| Z165 | `POST /group/inviteinfo` | núcleo adaptável: `POST group-info-from-invite-link` | núcleo adaptável: `POST group-info-from-invite-link` | J08 / P2 |
| Z166 | `POST /group/updateparticipants` | núcleo adaptável: `POST add-participant-group`; `POST remove-participant-group`; `POST promote-participant-group`; `POST demote-participant-group` | núcleo adaptável: `POST add-participant-group`; `POST remove-participant-group`; `POST promote-participant-group`; `POST demote-participant-group` | J08 / P2 |
| Z167 | `GET /admin/global/stats` | parcial: `GET /api/dashboard/stats`; `GET /metrics` | parcial: `GET /metrics` | J14 / P1 |
| Z168 | `POST /admin/global/test` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z169 | `GET /admin/global/config` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z170 | `POST /admin/global/config/reload` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z171 | `GET /admin/global/s3-retention/stats` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z172 | `POST /admin/global/s3-retention/trigger` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z173 | `POST /admin/global/event/test` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z174 | `GET /admin/events/stream` | parcial: `GET events/ws` | sem equivalente exposto | J14 / P1 |
| Z175 | `GET /admin/events/history` | parcial: `GET events` | sem equivalente exposto | J14 / P1 |
| Z176 | `GET /admin/events/stats` | parcial: `GET /metrics` | sem equivalente exposto | J14 / P1 |
| Z177 | `POST /admin/events/prune` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z178 | `POST /admin/events/session` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z179 | `GET /admin/dlq` | somente relacionado: `GET events` | sem equivalente exposto | J14 / P1 |
| Z180 | `DELETE /admin/dlq` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z181 | `GET /admin/dlq/stats` | parcial: `GET /metrics` | sem equivalente exposto | J14 / P1 |
| Z182 | `GET /admin/dlq/{id}` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z183 | `DELETE /admin/dlq/{id}` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z184 | `POST /admin/dlq/replay` | parcial: `POST events/replay` | sem equivalente exposto | J14 / P1 |
| Z185 | `POST /admin/dlq/prune` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z186 | `GET /admin/archive` | parcial: `GET events` | sem equivalente exposto | J14 / P1 |
| Z187 | `DELETE /admin/archive` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z188 | `GET /admin/archive/stats` | parcial: `GET /metrics` | sem equivalente exposto | J14 / P1 |
| Z189 | `GET /admin/archive/{id}` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z190 | `DELETE /admin/archive/{id}` | sem equivalente exposto | sem equivalente exposto | J14 / P1 |
| Z191 | `POST /admin/archive/replay` | parcial: `POST events/replay` | sem equivalente exposto | J14 / P1 |
| Z192 | `POST /status/send/text` | núcleo adaptável: `POST send-status-text` | núcleo adaptável: `POST send-text-storie` | J17 / P2 |
| Z193 | `POST /status/send/image` | núcleo adaptável: `POST send-status-media` | núcleo adaptável: `POST send-image-storie` | J17 / P2 |
| Z194 | `POST /status/send/video` | parcial: `POST send-status-media` | núcleo adaptável: `POST send-video-storie` | J17 / P2 |
| Z195 | `POST /status/send/audio` | sem equivalente exposto | sem equivalente exposto | J17 / P2 |
| Z196 | `POST /call/reject/send` | núcleo adaptável: `POST reject-call` | núcleo adaptável: `POST reject-call` | J13 / P4 |
| Z197 | `POST /call/reject` | núcleo adaptável: `POST reject-call` | núcleo adaptável: `POST reject-call` | J13 / P4 |
| Z198 | `POST /call/accept` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z199 | `POST /call/preaccept` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z200 | `POST /call/terminate` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z201 | `POST /call/initiate` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z202 | `GET /call/config` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z203 | `PUT /call/config` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z204 | `GET /call/status` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z205 | `POST /call/dial` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z206 | `POST /call/answer` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z207 | `POST /call/hangup` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z208 | `POST /call/play` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z209 | `POST /call/record/start` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z210 | `POST /call/record/stop` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z211 | `POST /call/group/dial` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z212 | `POST /call/group/dial-by-id` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z213 | `POST /call/participants/add` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z214 | `POST /call/participants/ring` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z215 | `POST /call/link/create` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z216 | `GET /call/link/preview` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z217 | `POST /call/link/join` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z218 | `POST /call/waiting-room/admit` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z219 | `POST /call/waiting-room/deny` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z220 | `POST /call/waiting-room/approval` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z221 | `POST /call/reaction` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z222 | `POST /call/hand` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z223 | `POST /call/video` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z224 | `POST /call/screenshare` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z225 | `GET /call/{call_id}/stream` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z226 | `GET /call/{call_id}/video/stream` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z227 | `GET /call/{call_id}/video/state` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z228 | `GET /call/{call_id}/video/stats` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z229 | `GET /call/{call_id}/events` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z230 | `GET /call/{call_id}/group/state` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z231 | `GET /call/{call_id}/waiting-room` | sem equivalente exposto | sem equivalente exposto | J13 / P4 |
| Z232 | `POST /chat/reply` | parcial: `POST send-reply` | núcleo adaptável: `POST send-reply` | J03 / P1 |
| Z233 | `POST /community/link` | núcleo adaptável: `POST link-community-group` | núcleo adaptável: `POST add-community-subgroup` | J09 / P3 |
| Z234 | `POST /community/unlink` | núcleo adaptável: `POST unlink-community-group` | núcleo adaptável: `POST remove-community-subgroup` | J09 / P3 |
| Z235 | `GET /community/linked` | parcial: `POST community-info` | sem equivalente exposto | J09 / P3 |
| Z236 | `POST /community/linked-groups` | parcial: `POST community-info` | sem equivalente exposto | J09 / P3 |
| Z237 | `POST /community/sub-groups` | parcial: `POST community-info` | sem equivalente exposto | J09 / P3 |
| Z238 | `POST /community/requests` | sem equivalente exposto | sem equivalente exposto | J09 / P3 |
| Z239 | `POST /community/requests/update` | sem equivalente exposto | sem equivalente exposto | J09 / P3 |
| Z240 | `POST /community/create` | núcleo adaptável: `POST create-community` | núcleo adaptável: `POST create-community` | J09 / P3 |
| Z241 | `POST /community/announcement` | composição: `POST send-message` | composição: `POST send-message` | J09 / P3 |
| Z242 | `POST /business/resolve-link` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z243 | `POST /business/contact-qr` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z244 | `GET /business/bots` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z245 | `POST /business/bot-profiles` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z246 | `POST /business/order/send` | sem equivalente exposto | parcial: `POST send-order-message` | J12 / P3 |
| Z247 | `GET /business/sticker-pack/{packId}` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z248 | `GET /business/order/{orderID}` | sem equivalente exposto | somente relacionado: `GET get-order-by-messageId/{messageId}` | J12 / P3 |
| Z249 | `GET /business/catalog/{businessJID}` | sem equivalente exposto | parcial: `GET get-products` | J11 / P3 |
| Z250 | `GET /business/products/{businessJID}` | sem equivalente exposto | parcial: `GET get-products` | J11 / P3 |
| Z251 | `GET /business/collections/{businessJID}` | sem equivalente exposto | parcial: `GET get-collections` | J11 / P3 |
| Z252 | `POST /business/catalog/send` | composição: `POST send-message` | parcial: `POST send-link-catalog` | J11 / P3 |
| Z253 | `POST /business/product/send` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z254 | `POST /business/product-list/send` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z255 | `POST /business/shop/send` | composição: `POST send-message` | parcial: `POST send-link-catalog` | J11 / P3 |
| Z256 | `POST /device/users` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z257 | `POST /device/users/context` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z258 | `GET /device/linked` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z259 | `POST /device/platform` | sem equivalente exposto | somente relacionado: `GET get-platform-from-message/{messageId}` | J07 / P2 |
| Z260 | `GET /privacy/settings` | núcleo adaptável: `GET privacy-settings` | sem equivalente exposto | J17 / P2 |
| Z261 | `POST /privacy/settings` | parcial: `POST privacy-settings` | sem equivalente exposto | J17 / P2 |
| Z262 | `POST /privacy/disappearing-timer` | parcial: `POST group-property` | núcleo adaptável: `POST temporary-messages` | J17 / P2 |
| Z263 | `GET /privacy/blocklist` | núcleo adaptável: `GET blocklist` | núcleo adaptável: `GET blocklist` | J17 / P2 |
| Z264 | `POST /privacy/blocklist` | núcleo adaptável: `POST block-contact`; `POST unblock-contact` | núcleo adaptável: `POST block-contact`; `POST unblock-contact` | J17 / P2 |
| Z265 | `GET /privacy/status` | sem equivalente exposto | sem equivalente exposto | J17 / P2 |
| Z266 | `PUT /business/cart` | sem equivalente exposto | núcleo adaptável: `POST set-cart-enabled` | J11 / P3 |
| Z267 | `POST /business/catalog` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z268 | `POST /business/catalog/{businessJID}/products/query` | sem equivalente exposto | parcial: `GET get-products` | J11 / P3 |
| Z269 | `GET /business/catalog/{businessJID}/products/{productID}` | sem equivalente exposto | parcial: `GET get-product-by-id` | J11 / P3 |
| Z270 | `POST /business/collections` | sem equivalente exposto | núcleo adaptável: `POST create-collection` | J11 / P3 |
| Z271 | `POST /business/collections/delete` | sem equivalente exposto | núcleo adaptável: `POST del-collection` | J11 / P3 |
| Z272 | `POST /business/collections/reorder` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z273 | `GET /business/collections/{businessJID}/{collectionID}` | sem equivalente exposto | composição: `GET get-collections` | J11 / P3 |
| Z274 | `PUT /business/collections/{collectionID}` | sem equivalente exposto | núcleo adaptável: `POST edit-collection` | J11 / P3 |
| Z275 | `POST /business/collections/{collectionID}/appeal` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z276 | `POST /business/cover` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z277 | `DELETE /business/cover/{coverID}` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z278 | `POST /business/eligibility` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z279 | `GET /business/linked-accounts` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z280 | `GET /business/merchant-compliance` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z281 | `PUT /business/merchant-compliance` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z282 | `POST /business/products` | sem equivalente exposto | núcleo adaptável: `POST add-product` | J11 / P3 |
| Z283 | `POST /business/products/delete` | sem equivalente exposto | núcleo adaptável: `POST del-products` | J11 / P3 |
| Z284 | `POST /business/products/image` | sem equivalente exposto | parcial: `POST add-product-image`; `POST change-product-image` | J11 / P3 |
| Z285 | `PUT /business/products/{productID}` | sem equivalente exposto | núcleo adaptável: `POST edit-product` | J11 / P3 |
| Z286 | `POST /business/products/{productID}/appeal` | sem equivalente exposto | sem equivalente exposto | J11 / P3 |
| Z287 | `PUT /business/products/{productID}/visibility` | sem equivalente exposto | núcleo adaptável: `POST set-product-visibility` | J11 / P3 |
| Z288 | `PUT /business/profile` | sem equivalente exposto | parcial: `POST edit-business-profile` | J11 / P3 |
| Z289 | `POST /chat/retry/message` | sem equivalente exposto | sem equivalente exposto | J03 / P1 |
| Z290 | `DELETE /newsletter/{newsletterID}` | sem equivalente exposto | núcleo adaptável: `DELETE newsletter/{id}` | J10 / P2 |
| Z291 | `POST /user/contact-link/resolve` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z292 | `GET /user/identity-code/{jid}` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z293 | `POST /user/username/resolve` | sem equivalente exposto | sem equivalente exposto | J07 / P2 |
| Z294 | `POST /group/historysharemode` | sem equivalente exposto | sem equivalente exposto | J08 / P2 |
| Z295 | `GET /business/catalog/own` | sem equivalente exposto | parcial: `GET get-products` | J11 / P3 |
| Z296 | `POST /business/order/accept` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
| Z297 | `POST /business/order/charge/send` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
| Z298 | `POST /business/order/payment/status/send` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
| Z299 | `POST /business/order/status/send` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
| Z300 | `POST /business/order/reject` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
| Z301 | `POST /group/history/send` | sem equivalente exposto | sem equivalente exposto | J08 / P2 |
| Z302 | `GET /business/order/state` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
| Z303 | `GET /business/order/state/{messageID}` | sem equivalente exposto | sem equivalente exposto | J12 / P3 |
