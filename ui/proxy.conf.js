const BETA_HOST = 'https://beta-clients1.shikshakul.com';

function stripCookieRestrictions(proxy) {
  proxy.on('proxyRes', (proxyRes) => {
    const cookies = proxyRes.headers['set-cookie'];
    if (cookies) {
      proxyRes.headers['set-cookie'] = cookies.map((cookie) =>
        cookie
          .replace(/;\s*Domain=[^;]*/gi, '')
          .replace(/;\s*Secure/gi, '')
      );
    }
  });
}

module.exports = {
  '/api/iam': {
    target: BETA_HOST,
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api\/iam/, '/iam'),
    configure: stripCookieRestrictions,
  },
  '/api/v1': {
    target: BETA_HOST,
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api\/v1/, '/acad'),
    configure: stripCookieRestrictions,
  },
  '/api/docs': {
    target: BETA_HOST,
    changeOrigin: true,
    rewrite: (path) => path.replace(/^\/api\/docs/, '/utity/docs'),
    configure: stripCookieRestrictions,
  },
};
