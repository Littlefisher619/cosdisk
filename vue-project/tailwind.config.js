module.exports = {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'media', // or 'media' or 'class'
  theme: {
    extend: {
      backgroundImage: {
        '2-texture': "url(https://cdn.jsdelivr.net/gh/MarleneJiang/ImgHosting/img/202109041905856.png)",
       }
    },
  },
  variants: {
  },
  plugins: [],
}