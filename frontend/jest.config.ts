import type { Config } from 'jest'
import nextJest from 'next/jest'

const createJestConfig = nextJest({
  // next.config.jsとテスト環境用の.envファイルを読み込むためのパスを提供
  dir: './',
})

// Jestのカスタム設定
const config: Config = {
  coverageProvider: 'v8',
  testEnvironment: 'jsdom',
  setupFilesAfterEnv: ['<rootDir>/jest.setup.ts'],
  moduleNameMapper: {
    // パスエイリアスの設定（tsconfig.jsonと一致させる）
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  testMatch: [
    '**/__tests__/**/*.[jt]s?(x)',
    '**/?(*.)+(spec|test).[jt]s?(x)',
  ],
  collectCoverageFrom: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.d.ts',
    '!src/**/*.stories.{js,jsx,ts,tsx}',
    '!src/**/__tests__/**',
  ],
}

// createJestConfigはこのようにエクスポートされ、next/jestが非同期でNext.jsの設定を読み込めるようにします
export default createJestConfig(config)
