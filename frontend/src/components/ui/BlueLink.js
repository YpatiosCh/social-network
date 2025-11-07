import Link from "next/link"

const BlueLinkButton = ({ where, what}) => {
    return (
        <Link
            href={where}
            className="link-primary"
        >
            {what}
        </Link>
    )
}

export default BlueLinkButton