const WhiteGrayButton = ({onClick=null, something=null, text=null}) => {
    return (
        <button
            onClick={onClick}
            className="btn-secondary"
        >
            {something}
            <span className="btn-text">{text}</span>
        </button>
    )
}

export default WhiteGrayButton