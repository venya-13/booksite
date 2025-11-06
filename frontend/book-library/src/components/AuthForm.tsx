import React, { useState } from 'react';
import { Box, TextField, Button, Typography, Stack, Divider, Paper } from '@mui/material';

interface AuthFormProps {
  type: 'login' | 'register';
  onSubmit: (email: string, password: string) => void;
  onGoogleAuth: () => void;
  onlyGoogle?: boolean;
}

const AuthForm: React.FC<AuthFormProps> = ({ type, onSubmit, onGoogleAuth, onlyGoogle }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(email, password);
  };

  return (
    <Box sx={{ maxWidth: 440, mx: 'auto', mt: 6, px: 2 }}>
      <Paper elevation={0} sx={{ p: 4, borderRadius: 4 }}>
        <Box component="form" onSubmit={handleSubmit}>
          <Stack spacing={2}>
            <Typography variant="h5" align="center" sx={{ fontWeight: 800 }}>
              {type === 'login' ? 'Welcome back' : 'Create your account'}
            </Typography>
            {!onlyGoogle && (
              <>
                <TextField
                  label="Email"
                  type="email"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  required
                />
                <TextField
                  label="Password"
                  type="password"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  required
                />
                <Button type="submit" variant="contained" color="primary" fullWidth>
                  {type === 'login' ? 'Login' : 'Register'}
                </Button>
                <Divider>or</Divider>
              </>
            )}
            <Button variant="outlined" color="secondary" fullWidth onClick={onGoogleAuth}>
              Continue with Google
            </Button>
          </Stack>
        </Box>
      </Paper>
    </Box>
  );
};

export default AuthForm; 