# RPC

This is project is **PRE-ALPHA** work-in-progress. Use at your own risk.

RPC is an attempt at building a simple cross-language RPC framework that tries
to provide gRPC-style code generation but without all the over-engineered
features that are only useful at Google scale. For example, we are never even
going to think about rolling a custom HTTP2 server implementation or requiring
developers to install 4 different support libraries just to compile a
basic "Hello, World" program.

Put another word, this is an RPC implementation that is designed for
startup-scale work favoring simplicity and ergonomics over extreme performance.

### Goals

While I do not think these goals will be unchallenged for long, I still believe
it is worth trying. In order to make the most easy-to-use developer-centric RPC
package, we need to solve many hard problems.

In order to guide the architectural decisions that we may face, I'd like to
propose:

1. **Developer Ergonomics before Runtime Performance.**  
   You should never need to think about whether or not your cloud provider
   supports a custom HTTP/2 extension in order to use RPC. It should just work.
   
   RPC uses the target platform's HTTP and JSON implementation whenever
   possible. There is no gain to be made from re-implementing web technologies
   that already works.
   
2. **Fat Binary and Monolithic Designs until really, really unavoidable.**  
   Plugin-based architecture is an interesting engineering problem. But they
   do not contribute to end-user productivity. Whenever there is a design or
   architectural choice to be made, we will choose the one which make the user
   more productive first, always.
   
   RPC is a single binary that will run the same way anywhere. You should be
   able to install a single binary, execute it with the same spec file, and get
   the same output.
   
3. **Predictable Fallbacks over Precision Semantics**  
   There is no need to have strongly-typed enums when a simple string field
   suffice in a language without them.
   
   There is no need to have perfect target language semantic (like `nil` vs `""`)
   when we can just as easily side-step it in the RPC spec (have no `nil` value)
   or use other simple predictable feature in the target language to get by.
   
   RPC will use the simplest possible fallback as often as needed. Even if it's
   not the prettiest or purest. If there is a potential for bikeshedding to take
   place, RPC will choose the choice that is the most boring.

# Usage

Install:

```sh
$ go get github.com/chakrit/rpc
```

Run:

```sh
$ rpc -gen go -out /api todo.rpc
```

* `-gen (lang)` – Currently supports Go and Elm, for now.
* `-out (folder)` - Outputs to specified folder.
* `todo.rpc` - The RPC spec file.

Develop:

```sh
$ # make edits
$ go get -v -u github.com/chakrit/smoke
$ ./test.sh
```

# Spec File

Syntax is a small language with braces.

```
option go_import "github.com/chakrit/todo/rpc"
option go_package "rpc"
option elm_module "Rpc"

namespace todo {
  enum State {
    Pending
    Completed
  }

  type TodoItem {
    string text
    State state
  }

  rpc ListItems() list<TodoItem>
  rpc AddItem(TodoItem) TodoItem
}
```

* `option __name__ __value__` - Sets target-specific option.
* `namespace __name__ { }` - Defines a scope.
* `type __name__ { }` - Defines an object type (or class or message).
* `enum __name__ { }` - Defines an enumeration of values.
* `rpc __name__ ( __args__ ) __return_args__` - Defines an RPC call.

Basic types:

| Name   | What's generated
| :--:   | :--
| string | Strings
| bool   | Booleans
| int    | Default integer type. 
| long   | 64-bit variant integer type, if available.
| float  | Default floating-point type.
| double | 64-bit variant floating-point type, if available.
| list   | Arrays or native list type.
| map    | Dictionaries or hashes.
| time   | Native time type, or same as `double` representing unix seconds.
| data   | Raw data buffers.

# LICENSE

See the LICENSE file included with the repository.

