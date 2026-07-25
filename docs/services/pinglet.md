# Pinglet

Upstream docs: https://pinglet.co.uk

Pinglet delivers notifications to a topic within a namespace via a simple
authenticated POST. Topics are auto-created on the first publish.

## URL Format

*pinglet://__`token`__@__`host`__/__`namespace`__/__`topic`__*

Use the `pinglet://` scheme (or set `scheme=http`) for an insecure endpoint; the
default `scheme=https` posts over TLS.

Badges (rendered as pills on the notification card, up to 3) and metadata (shown
on the detail sheet) are passed as comma-separated key/value pairs:

*pinglet://__`token`__@__`host`__/__`namespace`__/__`topic`__?priority=urgent&badges=Host:web-1,CPU:95%25&data=region:eu-west*

--8<-- "docs/services/pinglet/config.md"
