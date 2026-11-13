# Backend structure

Especially look at:

* common/joiner/
* common/starter/
* common/server/server_http/server_http_jschmhr
* common/auth/auth_server_http/endpoints.go
* apps/demo/demo_settings/starter.go
* apps/demo/demo.go


## The common application flow

* the full components list are created in apps/demo/demo_settings/components.go;
* each component is described using common/starter.Operator interface;
* apps/demo/demo.go reads a config, creates a logger and passes each component through "starter"
  (using, may be, some specific run-time options for it - so, the same component can be used twice
  in the same app);
* starting component gets these config and logger, creates own programming instances and passes them
  into the system.


## Joining instances between different components

We use import between components to pass only "static items" - values or functions - but not interface
instances.

To pass interface instance into the system we need to call joiner.Operator.Join() method (look at,
for example, into apps/demo/demo_settings/starter.go). To identify such a common instance
we use a unique key typed as joiner.InterfaceKey.

To get interface instance from the system (it should be previously stored using joiner.Operator.Join())
we need to call joiner.Operator.Interface() method (look at, for example, into
apps/demo/demo_settings/starter.go too). To identify such a common instance
we use an unique key typed as joiner.InterfaceKey.

Usually we can use predefined keys for each interface instances that are accessible via static import
from, for example, common/auth/keys.go. But if we need to use some interface twice or if we need some
specific key for it this key should be defined in call parameters for corresponding starters (the same
key for starter passing the instance and for starter that will get it).

