package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	Puerto            string
	TiempoEspera      time.Duration
	JWTSecret         string
	DemoAdminPassword string
	DemoUserPassword  string
}

func Cargar() (Config, error) {
	_ = godotenv.Load()

	var c Config

	c.DatabaseURL = os.Getenv("DATABASE_URL")
	if c.DatabaseURL == "" {
		return c, fmt.Errorf("falta DATABASE_URL")
	}

	c.Puerto = os.Getenv("PUERTO")
	if c.Puerto == "" {
		c.Puerto = "8080"
	}

	segundos := os.Getenv("TIEMPO_ESPERA_SEGUNDOS")
	if segundos == "" {
		segundos = "5"
	}

	n, err := strconv.Atoi(segundos)
	if err != nil || n <= 0 {
		return c, fmt.Errorf(
			"TIEMPO_ESPERA_SEGUNDOS debe ser un entero mayor que cero",
		)
	}

	c.TiempoEspera = time.Duration(n) * time.Second

	c.JWTSecret = os.Getenv("JWT_SECRET")
	if len(c.JWTSecret) < 32 {
		return c, fmt.Errorf("JWT_SECRET debe tener al menos 32 caracteres")
	}

	c.DemoAdminPassword = os.Getenv("DEMO_ADMIN_PASSWORD")
	if c.DemoAdminPassword == "" {
		return c, fmt.Errorf("falta DEMO_ADMIN_PASSWORD")
	}

	c.DemoUserPassword = os.Getenv("DEMO_USER_PASSWORD")
	if c.DemoUserPassword == "" {
		return c, fmt.Errorf("falta DEMO_USER_PASSWORD")
	}

	return c, nil
}
