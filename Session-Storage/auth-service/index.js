import express from 'express';
import session from 'express-session';
import Redis from 'ioredis';
import { RedisStore } from 'connect-redis';

const app = express();
app.use(express.json());

const redis = new Redis({
  host: process.env.REDIS_HOST || 'redis',
  port: process.env.REDIS_PORT || 6379,
});

redis.on('connect', () => {
  console.log('Connected to Redis successfully');
});

redis.on('error', (err) => {
  console.error(`Redis Error: ${err.message}`);
});

app.use(
  session({
    store: new RedisStore({ client: redis, }),
    secret: 'supersecret',
    resave: false,
    saveUninitialized: false,
    cookie: { secure: false, maxAge: 60000 },
  })
);

app.post('/login', (req, res) => {
  req.session.user = { id: 1, name: 'Alice' };
  res.send('Logged in and session saved');
});

app.get('/me', (req, res) => {
  res.send(req.session.user || {});
});

app.listen(3001, () => console.log('Auth service on port 3001'));
