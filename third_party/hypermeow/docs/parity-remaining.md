# Critérios das próximas entregas

O [plano](../PLAN.md) registra os checks. Use a [matriz congelada](parity-contract-matrix.md) como referência da ZuckZapGo v2.4.2 e o [progresso por contrato](parity-progress.md) para mudanças posteriores. Abaixo, cada entrega tem um resultado observável. Uma rota que responde 200 não fecha uma jornada.

| Entrega em aberto | Próxima implementação verificável | Evidência para marcar concluído |
| --- | --- | --- |
| Chatwoot | O provisionamento de um contato e sua conversa passou com a API real. Implementar anexos nos dois sentidos com limites e validação de destino. | Texto e anexo nos dois sentidos em inbox de teste. Uma resposta perdida não duplica conversa nem mensagem; mudança de conta/inbox, evento repetido, nota privada e eco não disparam envio. |
| Arquivamento automático | A reconciliação de upload explícito já passou com MinIO. Implementar opt-in por sessão, fila persistida com limite e lease, descoberta de anexos elegíveis e download limitado para o bucket privado. Não retroagir ao histórico, mídia excluída, temporária ou de visualização única. | S3 de teste: falha antes/depois de gravar, reinício, reconciliação por hash, limpeza de órfãos, quota e expiração. Uma segunda sessão não consegue listar, baixar nem apagar o arquivo. |
| Contatos e comunidades | Conferir origem de PN/LID, autorização de administrador, cursores e eventos em operações da matriz ainda parciais. | Fixtures do SDK/HTTP por comportamento; duas contas para confirmar efeitos de proprietário e participante quando o teste for retomado. |
| Comércio e pagamentos | Validar em catálogo Business o efeito das alterações de produtos, coleções, visibilidade e carrinho já cobertas em testes isolados. Tratar pedidos e estados financeiros com controle por papel e reconciliação. | Conta Business de teste com produto e coleção controlados. Preservar moeda e unidade. Mensagem de cobrança ou status do protocolo não comprova pagamento; nenhuma transação financeira real faz parte do teste automático. |
| Flows e interativos | Contratos distintos para botões, listas, imagem/vídeo, carrossel e Flow; identificar quais operações exigem Flow publicado e capacidade da conta. | Fixtures de protocolo por tipo. Recebimento, interação e retorno em Android, iPhone e Web permanecem adiados; o experimento não representa aceite nesses clientes. |
| Chamadas | Separar sinalização já observável de transmissão de áudio/vídeo. Documentar ciclo da chamada, autorização por sessão e limites antes de expor streams. | Duas contas de teste, áudio bidirecional, renegociação, perda de rede e encerramento. Medir liberação de CPU, memória e conexões. Rejeitar uma chamada não comprova suporte a chamadas de grupo. |
| Carga e restauração | Complementar o benchmark de consulta com um perfil acordado de conexões, mídia, falhas do broker, backlog, reconexão e restauração de backup. | Publicar ambiente, amostra, percentis, uso de CPU/memória/disco e perdas/duplicatas. As medições já feitas não são um SLA. |
| Monitoramento | Adiado por decisão do usuário: nenhum serviço nem destino de notificação configurado. | Quando retomado, entregar disparo e recuperação ao destino identificado. Exemplos de regras não equivalem a alertas ativos. |

Para o Chatwoot, os contratos oficiais distinguem [contatos](https://developers.chatwoot.com/api-reference/contacts/create-contact), [mensagens de conversa](https://developers.chatwoot.com/api-reference/messages/create-new-message) e [webhooks](https://developers.chatwoot.com/api-reference/webhooks/add-a-webhook). A implementação deve demonstrar como reconcilia cada recurso; um webhook genérico ou o workflow n8n não substituem o teste bidirecional.

O ensaio de queda de processo e backlog foi incorporado em [server-go #28](https://github.com/wppconnect-team/wppconnect-server-go/pull/28). Ele mede recuperação da fila SQLite após queda do processo. No [ensaio complementar com Redis](https://github.com/wppconnect-team/wppconnect-server-go/pull/33), quatro workers processaram 32 eventos enquanto o stream rejeitava as gravações. Depois de reabrir o banco e restaurar o stream, o replay preservou IDs, sessões e entregas confirmadas no webhook. Queda de rede e metas de capacidade continuam em aberto. O adaptador publicado recebe respostas públicas de atendentes e vincula mensagens diretas a uma conversa. O [teste com a aplicação Chatwoot](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42911977000) exercita a API de uma inbox em PostgreSQL e a assinatura de um callback do atendente. O lado WhatsApp continua simulado nesse teste. Para o sentido inverso, o comando de mensagem aceita JSON para texto e multipart para anexos na [API oficial do Chatwoot](https://developers.chatwoot.com/api-reference/messages/create-new-message); isso por si só não prova deduplicação.

## Catálogo da sessão pareada

O operador confirmou que a sessão pareada usa WhatsApp comum e não tem catálogo. Nas duas consultas, o Manager devolveu HTTP 502 e o SDK informou GraphQL 400. Não houve escrita. O aceite no WhatsApp da leitura, criação, edição e exclusão de produtos, das alterações de coleções, de visibilidade e de carrinho fica pendente para uma conta WhatsApp Business com catálogo de teste. A sessão pessoal continua disponível para os testes de mensagens. [Registro da sessão](catalog-live-account-acceptance.json).

## Execução da referência

A imagem AMD64 da matriz foi iniciada duas vezes em uma rede Docker interna com um PostgreSQL temporário pronto. Em ambas as tentativas, o processo encerrou com código 1 e informou ausência de `LICENSE_KEY`. O operador confirmou que não dispõe de uma instalação licenciada ou chave de teste. A comparação comportamental por HTTP continua pendente. O CI dos adaptadores Go e a comparação dos dois OpenAPIs podem ser reproduzidos separadamente. [Detalhes do ensaio](reference-runtime-acceptance.json).

## Comparação atualizada do servidor Node com o Go

O Swagger da revisão de referência do Node, `37b7bc1a662b`, descreve 158 operações. O OpenAPI do Go em `1293a52d5bfd` descreve 170 operações em 165 caminhos. A interseção por **método e caminho literais** tem 87 operações. Restam 83 combinações exclusivas do Go e 71 exclusivas do Node. Essa contagem inclui capacidades experimentais ou indisponíveis quando presentes no contrato. [Swagger do Node](https://github.com/wppconnect-team/wppconnect-server/blob/37b7bc1a662b5262b8276597701f38c4dbcad58a/src/swagger.json) e [OpenAPI do Go](https://github.com/wppconnect-team/wppconnect-server-go/blob/1293a52d5bfd772ae8d8e0d53116465e295bd2fd/docs/openapi.json).

Para reproduzir a contagem com os dois repositórios locais, execute `python docs/compare-server-contracts.py CAMINHO_GO CAMINHO_NODE`. O script lê os blobs Git fixados pelo plano, lista as contagens por jornada e identifica os candidatos abaixo.

Os dez diferenciais inicialmente identificados entre o núcleo adaptável do Node e o Go têm contratos Go nativos testados. Z136, Z266, Z270, Z274 e Z287 foram registrados nas fases anteriores. Os cinco abaixo passaram na suíte isolada do servidor. [Guia dos cinco comandos](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/canal-e-catalogo.pt-BR.md).

| Referência | Operação | Limite do aceite |
| --- | --- | --- |
| Z290 | Remover um canal próprio com o mesmo ID e confirmação `DELETED`. | Exclusão real adiada. A suíte verifica recusa por papel, cliente incompatível, resposta incerta e repetição. |
| Z271 | Excluir de 1 a 100 coleções no catálogo da sessão. | A confirmação cobre o lote. A verificação individual exige reler as páginas do catálogo. |
| Z282, Z285 | Criar e editar a ficha do produto com referências HTTPS de mídia, moeda e preços inteiros em milésimos. | O ID, a visibilidade explícita em edição e a resposta de revisão passaram com dados sintéticos. Falta observar o estado no vendedor e no comprador. |
| Z283 | Excluir de 1 a 100 produtos. | O servidor informa contagem total ou parcial. A identificação dos itens remanescentes depende da leitura do catálogo. |

[Handler de favoritos do Node](https://github.com/wppconnect-team/wppconnect-server/blob/37b7bc1a662b5262b8276597701f38c4dbcad58a/src/controller/deviceController.ts#L2125), [patch de favorito](https://github.com/wppconnect-team/hypermeow/blob/64b6d5ba94fa5610ffe6417cc19dfbaa6e3baffb/appstate/encode.go#L303) e [remoção de canal no SDK](https://github.com/wppconnect-team/hypermeow/blob/64b6d5ba94fa5610ffe6417cc19dfbaa6e3baffb/newsletter_delete.go).

## Cinco contratos em relação ao Swagger da referência

Comparação estática com `/app/static/api/swagger.json` na imagem v2.4.2. SHA256 do arquivo: `2c15e65a226a2a8bae832bdfcfcf3ecbb3b990eec1e71e29a4272a50a0c98d21`. O guia do Go documenta o método e o corpo que o cliente deve usar na Cloud.

| ID | Método e caminho da Zuck | Comando Go (POST) |
| --- | --- | --- |
| Z290 | `DELETE /newsletter/{newsletterID}` | `delete-newsletter` recebe `id` no JSON. |
| Z271 | `POST /business/collections/delete` | `del-collection` recebe `collectionIds`. A Zuck documenta `collection_ids`. |
| Z282 | `POST /business/products` | `add-product` recebe `product` e devolve 201. A Zuck documenta 200. |
| Z285 | `PUT /business/products/{productID}` | `edit-product` recebe `productId` e `product`. Exige visibilidade explícita ao substituir a ficha. |
| Z283 | `POST /business/products/delete` | `del-products` recebe `productIds`. A Zuck documenta `product_ids`. |

Os dois contratos usam strings inteiras em milésimos para preços e referências de mídia já hospedada. A especificação da Zuck oferece largura e altura opcionais da prévia. O Go solicita 100 × 100 ao SDK e devolve a projeção de produto do servidor com o estado de revisão. Isso evita persistir URLs temporárias de mídia no recibo. A edição substitui a ficha inteira. Mantenha no sistema de origem a descrição, as mídias, o preço promocional e os dados de conformidade que deseja preservar. A leitura resumida do Go não contém todos esses campos.

O Swagger da referência descreve as respostas de sucesso como objetos livres. Só uma execução licenciada pode confirmar seus corpos reais, os códigos de falha e os efeitos no WhatsApp. O [OpenAPI do Go](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/openapi.json) e o [guia do servidor](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/canal-e-catalogo.pt-BR.md) descrevem o contrato nativo verificado.
