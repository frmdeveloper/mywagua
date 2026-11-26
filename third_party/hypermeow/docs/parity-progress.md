# Progresso sobre a matriz Zuck v2.4.2

A [matriz de 303 operações](parity-contract-matrix.md) é o inventário da versão comparada. Esta tabela registra o que mudou depois da medição, sem alterar o baseline nem somar rotas como prova de equivalência. Os [contratos de canais](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/newsletters.pt-BR.md) mostram os payloads nativos e a interpretação dos códigos HTTP.

| ID | Comportamento implementado | Limite da evidência |
| --- | --- | --- |
| Z014 | Consultar metadados por convite resolvido no WhatsApp | Teste HTTP verifica a identidade do canal; não é equivalência de todos os campos. |
| Z017, Z019 | Seguir/deixar a partir do convite | O SDK recebe o JID resolvido; a resposta contém os metadados anteriores à alteração. |
| Z022 | Enviar/remover reação por ID de servidor | 202 significa gravação no transporte; sem confirmação do cliente remoto. |
| Z023 | Registrar visualização e esperar ACK | Fixture do SDK cobre ACK, rejeição, desconexão e cancelamento; repetição com chave não duplica o efeito no teste. |
| Z025 | Consultar contadores e atualizações com `since` e `after` | Janela limitada, sem filtro `before` e sem garantia de histórico completo. |
| Z026 | Aceitar aviso específico com consentimento explícito | Testado com cliente simulado. Nenhum termo foi aceito por conta real. |
| Z290 | Excluir canal após verificar plataforma, ID e papel de dono. Uma confirmação válida identifica o mesmo canal com estado `DELETED`. | Cliente simulado e recibo SQLite. Exclusão real adiada por escolha do operador. |
| Z027 | Enviar convite como administrador ou dono verificado | Verificação de papel e mensagem de protocolo testadas. Entrega, aceite e promoção não foram validados em duas contas. |

O [recibo de canary, pool e navegador](channels-acceptance-2026-11-10.json) não contém telefones, credenciais nem identificadores das sessões. O [plano](../PLAN.md) separa CI, publicação e aceite real. Os testes que exigem segunda conta ou Android/iPhone/Web foram adiados pelo usuário.

## Comércio, brokers e interativos

| Área | Evidência automatizada | O que permanece em aberto |
| --- | --- | --- |
| J11: Z249, Z250, Z251 e Z269 | [Catálogo da conta conectada](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/catalogo.pt-BR.md), sem substituir vendedor, preservando moeda, preço e visibilidade falsa. | Catálogo de outra empresa, pedidos, pagamentos e validação com vendedor real. |
| J11: Z266 | Configurar o carrinho do vendedor com `set-cart-enabled`. O SDK valida o booleano devolvido pelo WhatsApp. Os testes HTTP exercitam habilitar, desabilitar e recuperar o resultado após falha. | Alteração do carrinho em uma conta Business e confirmação da experiência do comprador. |
| J11: Z270, Z274 e Z287 | Criar e editar coleções e alterar a visibilidade de produtos. Escopo de escrita e consulta de resultado próprios, repetição idempotente, conflito e resposta remota incerta testados. | Moderação e visibilidade para o comprador em conta Business real. |
| J11: Z271, Z282, Z283 e Z285 | Excluir coleções, criar, editar e excluir produtos no catálogo da sessão. A suíte cobre preço e moeda, resposta de revisão, contagem parcial e repetição após falha. | Confirmação individual de coleções, revisão comercial e visibilidade para o comprador aguardam conta Business. |
| J15: Redis Streams e SQS | [Dois destinos opcionais](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/event-delivery-operations.md) com identidade do evento, recuperação e inspeção. Integração com Redis real e com ElasticMQ para SQS. | Recursos IAM/KMS na AWS, configuração CRUD por sessão e aceite operacional do consumidor. |
| J05: Z126 | Carrossel experimental de imagens e respostas rápidas, montado pelo SDK. Testes de JSON, upload, protobuf, autorização e replay. | Vídeo, outros controles, entrega e apresentação em Android/iPhone/Web. |

O teste [n8n 2.47.3](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42902244064/job/130130987174) importou e publicou o workflow, rejeitou uma credencial inválida, aceitou a credencial da fixture e retornou somente ID do evento, nome e sessão. Nenhum evento foi enviado a um n8n de cliente.

A [sobreposição legível por máquina](parity-progress.json) referencia os IDs da matriz, os métodos nativos, campos obrigatórios, códigos HTTP, idempotência, nomes dos testes e limites de cada evidência. Execute `python docs/check-parity-progress.py CAMINHO_DO_SERVER_GO` para verificar o contrato de um checkout. Com `--run-fixtures`, o validador executa os testes HTTP mapeados e exige um resultado aprovado para cada função. O CI verifica o contrato e executa essa suíte no SHA declarado em `serverRevision`. O aceite em dispositivo segue registrado à parte.

A [comparação das duas imagens](session-verification-benchmark.json) mede o endpoint de recuperação autenticada. O [canary da última imagem](session-cache-acceptance.json) e o [recibo da publicação do n8n](workflows-carousel-acceptance.json) distinguem as verificações na pool dos testes de protocolo. Os trabalhos restantes, critérios de aceite e adiamentos estão em [próximas entregas](parity-remaining.md).

## Chatwoot

O [adaptador de texto](https://github.com/wppconnect-team/wppconnect-server-go/blob/main/docs/chatwoot-text-relay.md) liga uma sessão Cloud a uma conversa existente em uma inbox de API. Para respostas públicas do atendente, o recibo aponta para a operação de envio da Cloud. Para texto direto recebido do WhatsApp, o identificador da mensagem permite consultar a gravação correspondente no Chatwoot após uma resposta incerta.

No [teste com Chatwoot 4.19.1](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42911977000/job/130164677043), o Chatwoot persistiu uma mensagem após duas entregas sintéticas do mesmo evento. O caminho inverso usou o callback assinado da aplicação e um servidor HTTP local para representar a Cloud. O [teste de provisionamento](https://github.com/wppconnect-team/wppconnect-server-go/actions/runs/42915535317) confirmou a recuperação dos IDs após interrupção da resposta de criação de contato, conversa e vínculo de inbox. Anexos e entrega a um telefone seguem nos [critérios de aceite](parity-remaining.md).

## Favoritos e sincronização de estado

Z136 usa `star-message`, `message-state-result` e `message-star-state`. O primeiro comando altera a preferência, o segundo recupera o recibo e o último consulta o booleano efetivamente observado. A sessão procura a mensagem e o remetente no próprio histórico. O teste cobre repetição, confirmação incerta, colisão de chave e isolamento entre sessões. [Contrato do servidor](https://github.com/wppconnect-team/wppconnect-server-go/blob/a7c46b0c3a1c127f6f3b9f969ff2a3b3bdcb520d/docs/favoritos.pt-BR.md).
