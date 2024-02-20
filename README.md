# Weather Service (Take Home)

## API Output

The API provides a JSON response that includes the current weather details for the specified latitude and longitude. The response contains the following fields:

- `description`: A summary of the current weather condition

- `temperature_category`: The temperature categorized in the following way...
  - `cold`: At or below 10°C (50°F).
  - `moderate`: Above 10°C (50°F) and up to 25°C (77°F).
  - `hot`: Above 25°C (77°F).

Here's an example of a typical API response:

```json
{
  "description": "light rain",
  "temperature_category": "moderate"
}
```

## Building the image

`docker build -t weather-app .`

## Running the service

`docker run -p 8080:8080 -e OPENWEATHER_API_KEY=(your api key) weather-app`

## Testing

`curl "http://localhost:8080/weather?lat=37&lon=-122"`

## Improvements

To make this service more production-ready, the following improvements could be made...

- **External Configuration**: Use external configuration files or environment variable management tools to configure the service dynamically. This includes server port, OpenWeather API base URL, version, and temperature category mappings.

- **Rate Limiting**: Implement rate limiting through middleware that tracks and limits requests per client over a defined time window.

- **Timeouts**: Set timeouts for API calls to avoid hanging requests by utilizing context with deadlines.

- **Structured Logging**: Adopt structured logging to improve log readability and parsing by log management systems.

- **Monitoring and Tracing**: Integrate monitoring and tracing tools to track health, performance metrics, and request traces.

- **TLS/SSL**: Secure all communications with TLS.
