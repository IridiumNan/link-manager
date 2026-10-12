# Link Manager

This is a single tool that manage system soft links conveniently.

The configuration file will be a single map link this

```json
{
    "dir_map": {
        "src_dir": "dst_dir"
    },
    "file_map": {
        "old_name": "new_name"
    }
}
```

> [!NOTE]
> Just use the absolute path on configuration file

## Quick Start

- run program (it will create a new configuration file)

- edit your configuration file

- use `lm run` to manage or `lm install` to install all

---

## Why

As there is always some link share the src dir and dst dir, so manage them as groups will be faster.

For example, you want to manage your service unit file & scripts utils across your machine.

So you just assemble them on ~/utils then create soft-links to `~/.config/systemd/user/` & `~/.local/bin`

And this project is powered by [buttletea](https://github.com/charmbracelet/bubbletea)
