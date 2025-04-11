import express from 'express';
import session from 'express-session';
import { RedisStore } from 'connect-redis';
import Redis from 'ioredis';

const redisClient = new Redis({ host: 'redis', port: 6379 });

redisClient.on('connect', () => {
  console.log('Connected to Redis from Profile Service');
});
redisClient.on('error', err => {
  console.error('Redis Error:', err);
});

const app = express();
app.use(express.json());

// Use session middleware
app.use(
  session({
    store: new RedisStore({ client: redisClient }),
    secret: 'supersecret',
    resave: false,
    saveUninitialized: false,
    cookie: {
      secure: false, // set true if behind HTTPS
      maxAge: 60000
    }
  })
);

// Sample route
app.get('/profile', (req, res) => {
  if (!req.session.user) {
    return res.status(401).send('Not authenticated');
  }
  res.send(`Welcome ${req.session.user.name}!`);
});

app.listen(3002, () => {
  console.log('Profile service running on port 3002');
});
