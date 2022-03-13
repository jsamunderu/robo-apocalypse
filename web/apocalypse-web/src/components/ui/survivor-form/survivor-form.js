import classes from './survivor-form.module.css';
const SurvivorForm = (props) => {
  return (
    <form className={classes.form}>
      <div className='form-fields'>
        <div>
          <label htmlFor="name">Name</label>
          <input type="text" id="name" />
        </div>
        <div>
          <label htmlFor="name">Age</label>
          <input type="text" id="age" />
        </div>
        <div>
          <label htmlFor="name">Gender</label>
          <input type="text" id="gener" />
        </div>
        <div>
          <label htmlFor="name">Id Number</label>
          <input type="text" id="idnumber" />
        </div>
      </div>
      <div className="form-actions">
        <button>Submit</button>
      </div>
    </form>
  );
};

export default SurvivorForm;
