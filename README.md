![Fatsecret_CLI_(unoffical)](https://github.com/user-attachments/assets/de9f0789-0157-4994-939c-a7a00cd58487)


Fatsecret CLI is a way to search for nutritional information about food in the terminal. 

Nutritional information can be searched directly from the fatsecret database. Search results can be stored in a favourites file for easy access. You can search through favourites and display their information and you can also delete favourites.


![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/LodewykRoux/fatsecret-cli/go.yml?branch=main&label=Build%20%26%20Test&color=fff)
[![Static Badge](https://img.shields.io/badge/go-v1.23.5-blue?logo=go&color=%2300ADD8)](https://go.dev/)
[![Static Badge](https://img.shields.io/badge/cobra-v1.8.1-orange?color=%23efb52d)](https://github.com/spf13/cobra)
[![Static Badge](https://img.shields.io/badge/bubble--tea-v1.3.0-purple?color=%236c50ff)](https://github.com/charmbracelet/bubbletea)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)



## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`ENCRYPTION_KEY` - 16 Digit long encryption key. 




## Run Locally

Clone the project

```bash
  git clone https://github.com/LodewykRoux/fatsecret-cli
```

Go to the project directory

```bash
  cd fatsecret-cli

```
## API Reference

#### Set the client secret and the client id (get them [here](https://platform.fatsecret.com/platform-api))

```bash
  go run main.go config
```

#### Get food item from fatSecret api

```bash
  go run main.go food --term="something"
```

![til](./demo/food.gif)

#### Favourites
##### List all favourites
```bash
  go run main.go favourite list
```

![til](./demo/favourites.gif)

##### List all favourites, select to delete
```bash
  go run main.go favourite delete
```

![til](./demo/favourites-delete.gif)

##### List all favourites with similar name
```bash
  go run main.go favourite search --term="something"
```

![til](./demo/favourites-search.gif)





## License

[MIT](https://choosealicense.com/licenses/mit/)

