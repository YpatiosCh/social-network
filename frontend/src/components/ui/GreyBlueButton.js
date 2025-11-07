const GrayBlueButton = ({onClick=null, something=null, text=null }) => {
    return (
        <button
            onClick={onClick}
            className="btn-primary"
        >
            {something}
            <span className="btn-text">{text}</span>
        </button>
    )
}

export default GrayBlueButton