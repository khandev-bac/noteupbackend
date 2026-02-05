package utils

var (
	PORT                     string = ":3000"
	FIREBASE_SERVICE         string = GetEnv("FIREBASE_URL", "internals/config/noteon_back.json")
	FIREBASE_PROJECT_ID      string = GetEnv("FIREBASE_PROJECT_ID", "")
	PRODUCTION_DB_URL        string = GetEnv("DB_URL", "")
	ACCESSTOKEN_KEY          string = GetEnv("ACCESS_TOKEN_KEY", "")
	REFRESHTOKEN_KEY         string = GetEnv("REFRESH_TOKEN_KEY", "")
	DEVELOPMENT_DATABASE_URL string = GetEnv("DB_DEV_URL", "")
	APP_ENV                  string = GetEnv("APP_ENV", "")
	GEMINI_API               string = GetEnv("GEMINI_API_KEY", "")
	REDIS_HOST               string = GetEnv("REDIS_HOST", "localhost")
	REDIS_TTL                string = GetEnv("REDIS_TTL", "")
	IMAGEKIT_PRIVATE_KEY     string = GetEnv("IMAGEKIT_PRIVATE_KEY", "")
	IMAGE_KIT_BASE_URL       string = GetEnv("IMAGE_KIT_BASE_URL", "")
	NOTECREATEDEVENT         string = "note.created"
	AUDIONOTESTOPIC          string = "audio_notes"
)
var MAX_FREE_AUDIO_FILE_LIMIT = 10 << 20
var MAX_PREMIUN_AUDIO_FILE_LIMIT = 80 << 20
