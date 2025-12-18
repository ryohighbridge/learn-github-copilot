describe('基本的なテスト', () => {
  it('Jestが正しく動作する', () => {
    expect(1 + 1).toBe(2)
  })

  it('文字列の比較ができる', () => {
    expect('hello').toBe('hello')
  })

  it('配列の要素を確認できる', () => {
    const arr = [1, 2, 3]
    expect(arr).toHaveLength(3)
    expect(arr).toContain(2)
  })

  it('オブジェクトのプロパティを確認できる', () => {
    const obj = { name: 'Test', value: 123 }
    expect(obj).toHaveProperty('name')
    expect(obj.value).toBe(123)
  })
})
