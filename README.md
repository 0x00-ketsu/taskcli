# Taskcli
 ```markdown
 ________________________________
< A terminal UI for manage tasks >
 --------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
 ```

[](https://user-images.githubusercontent.com/16932133/196026551-b9b51e25-a35f-4f4c-bb38-aff134277564.mp4) 

## 🌟 Features

- Mange tasks in terminal

- All operations with keyboard

- Support `markdown` syntax

If you want to preview markdown in terminal, you can try my another project:
[smooth](https://github.com/0x00-ketsu/smooth)


## 📥 Install

### Go install
```shell
go install github.com/0x00-ketsu/taskcli@latest
```

### From source

```shell
git clone https://github.com/0x00-ketsu/taskcli
cd taskcli
make build
```

## 🔭 Usage

Input `taskcli` in terminal then `Enter`

```shell
taskcli
```

```shell
taskcli -h
A terminal UI for manage tasks

Usage:
  taskcli [flags]

Flags:
  -c, --editor string    external editor for task detail panel (default "vim")
  -h, --help             help for taskcli
  -s, --storage string   taskcli data storage location (default "~/.taskcli/bolt.db")
```

## 🖮 Keyboard shortcuts

| Scope    | Shortcut    | Action    |
|---------------- | --------------- | --------------- |
| Global    | `q`    | Quit Application    |
| Global    | `Esc`    |   Step back  |
|  Filter Panel   | `j`    |  Move to next item   |
|  Filter Panel   | `k`    |  Move to previous item   |
|  Filter Panel   |   `g`  |  Go to first item   |
|  Filter Panel   |   `G`  |  Go to last item   |
|  Filter Panel   |   `Enter`  | Activate task    |
|  Search Panel   |   `Esc`  |  Back to Filter panel   |
|  Search Panel   |   `Tab`  | Change field    |
|  Task Panel   | `n`    | Create a new task    |
|  Task Panel   | `j`    | Move to next item    |
|  Task Panel   | `k`    | Move to previous item    |
|  Task Panel  |  `g`   |  Go to first item   |
|  Task Panel  |  `G`   |  Go to last item   |
|  Task Panel  |  `m`   |  Show menus   |
|  Task Detail Panel   |  `r`   | Rename task title    |
|  Task Detail Panel   |  `t`   | Set to today    |
|  Task Detail Panel   |  `+`   | Set to next day    |
|  Task Detail Panel   |  `-`   | Set to previous day    |
|  Task Detail Panel   |  `space`   | Toggle task status    |
|  Task Detail Panel   |  `i`   | Edit task content view    |
|  Task Detail Panel   |  `c`   |  Copy task content   |
|  Task Detail Panel   |  `j`   |  Move cursor down   |
|  Task Detail Panel   |  `k`   |  Move cursor up   |
|  Task Detail Panel   |  `h`   |  Move cursor left  |
|  Task Detail Panel   |  `l`   |  Move cursor right  |
|  Task Detail Panel   |  `Ctrl-d`   |  Scroll down  |
|  Task Detail Panel   |  `Ctrl-u`   |  Scroll up  |
|  Task Detail Panel   |  `v`   |  Edit task content with external editor  |



## 💡 Inspiration by

[geek-life](https://github.com/ajaxray/geek-life) by Anis uddin Ahmad

## 🔖 License

MIT
