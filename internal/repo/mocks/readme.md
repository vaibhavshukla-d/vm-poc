# Understanding `testify/mock` and How `.On()` Works

If you’re writing unit tests in Go and using **Testify** with **mockery**, the `.On()` method is your main tool to **mock method calls**. Here’s a simple breakdown of what it does, how it works, and some practical examples.  

---

## 1️⃣ What `.On()` Does

`.On()` is used to **tell a mock object how to behave** when a specific method is called.

**Syntax:**
```go
mockObj.On(methodName string, arguments ...interface{}) *mock.Call


2️⃣ How .On() Works Internally

When you call .On(), Testify creates a Call object and stores it inside the mock.

The mock stores your expected arguments and return values.

Later, when the method is actually called:

Testify searches for a matching Call based on method name and arguments.

If it finds a match, it runs any .Run() callback and returns the values from .Return().

.Once() or .Times(n) make sure calls are only used the expected number of times.

In other words: .On() is like scripting the mock’s behavior.