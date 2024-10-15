/** @type {import('next').NextConfig} */
const nextConfig = {
    distDir: 'dist',
    output: 'export',
    swcMinify: false,
    experimental: {
        forceSwcTransforms: true,
    },
}

module.exports = nextConfig
