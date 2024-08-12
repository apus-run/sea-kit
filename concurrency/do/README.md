<table>
<tr>
<th><code>do</code></th>
<th><code>stdlib</code></th>
</tr>
<tr>
<td>

```go
err := Do(
    func() error {
        return doThingA(),
    },
    func() error {
        return doThingB(),
    },
    func() error {
        return doThingC(),
    })
```

</td>
<td>

```go
var wg sync.WaitGroup
var errs []error
errChan := make(chan error)

wg.Add(3)
go func() {
    defer wg.Done()
    if err := doThingA(); err != nil {
        errChan <- err
    }
}()
go func() {
    defer wg.Done()
    if err := doThingB(); err != nil {
        errChan <- err
    }
}()
go func() {
    defer wg.Done()
    if err := doThingC(); err != nil {
        errChan <- err
    }
}()

go func() {
    wg.Wait()
    close(errChan)
}()

for err := range errChan {
    errs = append(errs, err)
}

err := errors.Join(errs...)
```

</td>
</tr>
</table>