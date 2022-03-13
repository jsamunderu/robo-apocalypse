import classes from './header.module.css';
import theMatrix from './the-matrix.png'
const Header = (props) => {
  return (
    <>
      <header className={classes.header}>
        <h1>Robot Apocalypse</h1>
      </header>
      <div className={classes['main-image']}>
        <img src={theMatrix} alt='The matrix' />
      </div>
    </>
  );
};

export default Header;
